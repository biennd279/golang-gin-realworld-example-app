package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/articles"
	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
)

type appContext struct {
	db *gorm.DB
	r  *gin.Engine
}

func (a *appContext) reset() {
	a.r = gin.Default()

	v1 := a.r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	articles.TagsAnonymousRegister(v1.Group("/tags"))

	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))

	articles.ArticlesRegister(v1.Group("/articles"))

	a.db = common.TestDBInit()
	migrateDB(a.db)
}

func migrateDB(db *gorm.DB) {
	users.AutoMigrate()
	db.AutoMigrate(&articles.ArticleModel{})
	db.AutoMigrate(&articles.TagModel{})
	db.AutoMigrate(&articles.FavoriteModel{})
	db.AutoMigrate(&articles.ArticleUserModel{})
	db.AutoMigrate(&articles.CommentModel{})
}

func (a *appContext) teardown() {
	common.TestDBFree(a.db)
}

type appCtxKey struct{}

type response struct {
	statusCode int
	body       string
	headers    http.Header
}

type rspCtxKey struct{}

func ConvertToString(model interface{}) string {
	bytes, _ := json.Marshal(model)
	return string(bytes)
}

func NewJSONRequest(method string, target string, param interface{}) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(ConvertToString(param)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return req
}
