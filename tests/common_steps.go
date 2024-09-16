package tests

import (
	"context"
	"encoding/json"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/articles"
	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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

func iHaveAValidEmailAndPasswordIs(ctx context.Context, email string, password string) (context.Context, error) {

	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	userModel := users.UserModel{
		Username:     "test",
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	users.SaveOne(&userModel)

	return ctx, nil
}

func iHaveAInvalidUsernameAndPasswordIsInvalidAndInvalid(ctx context.Context) (context.Context, error) {
	// do nothing
	return ctx, nil
}

func iLoginWithTheValidEmailAndPassword(ctx context.Context) (context.Context, error) {
	return ctx, godog.ErrPending
}

func iHaveAValidToken(ctx context.Context) (context.Context, error) {
	return ctx, godog.ErrPending
}

func SetupDefaultApplicationScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		app := &appContext{}
		app.reset()
		return context.WithValue(ctx, appCtxKey{}, app), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		app := ctx.Value(appCtxKey{}).(*appContext)
		app.teardown()
		return ctx, nil
	})

	ctx.Step(`^I have a invalid username and password is "([^"]*)" and "([^"]*)"$`, iHaveAInvalidUsernameAndPasswordIsInvalidAndInvalid)
	ctx.Step(`^I have a valid email and password is "([^"]*)" and "([^"]*)"$`, iHaveAValidEmailAndPasswordIs)
	ctx.Step(`^I login with the valid email and password$`, iLoginWithTheValidEmailAndPassword)
	ctx.Step(`^I have a valid token$`, iHaveAValidToken)
	ctx.Step(`^I am unauthenticated with invalid token$`, iAmUnauthenticatedWithInvalidToken)
}

func iAmAuthenticatedWithValidToken(ctx context.Context, stateToken string) (context.Context, error) {
	return ctx, nil
}
func iAmUnauthenticatedWithInvalidToken(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
