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

func newAppContext() *appContext {
	app := &appContext{}
	app.r = gin.Default()

	gin.SetMode(gin.TestMode)

	v1 := app.r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	v1.Use(users.AuthMiddleware(false))
	articles.ArticlesAnonymousRegister(v1.Group("/articles"))
	articles.TagsAnonymousRegister(v1.Group("/tags"))

	v1.Use(users.AuthMiddleware(true))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))

	articles.ArticlesRegister(v1.Group("/articles"))

	app.db = common.TestDBInit()
	app.db.LogMode(false)

	return app
}

func (a *appContext) reset() {
	a.db = common.TestDBInit()
	a.db.LogMode(false)
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

type tokenCtx struct {
	token string
}

type tokenCtxKey struct{}

type responseCtx struct {
	statusCode int
	body       string
	headers    http.Header
}

type rspCtxKey struct{}

type userLoginResponse struct {
	User struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Bio      string  `json:"bio"`
		Image    *string `json:"image"`
		Token    string  `json:"token"`
	} `json:"user"`
}

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

func CreateLoginRequest(email string, password string) *http.Request {
	param := gin.H{
		"user": gin.H{
			"email":    email,
			"password": password,
		},
	}
	return NewJSONRequest("POST", "/api/users/login", param)
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

func SetupDefaultApplicationScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		app := newAppContext()
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
	ctx.Step(`^I am unauthenticated with invalid token$`, iAmUnauthenticatedWithInvalidToken)
	ctx.Step(`^I am authenticated with valid token$`, iAmAuthenticatedWithValidToken)
}

func iAmAuthenticatedWithValidToken(ctx context.Context) (context.Context, error) {
	loginReq := CreateLoginRequest("test@gmail.com", "password")
	w := httptest.NewRecorder()
	app := ctx.Value(appCtxKey{}).(*appContext)
	app.r.ServeHTTP(w, loginReq)

	var loginResponse userLoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)

	if err != nil {
		return ctx, err
	}

	reqCtx := tokenCtx{
		token: loginResponse.User.Token,
	}

	return context.WithValue(ctx, tokenCtxKey{}, reqCtx), nil
}

func iAmUnauthenticatedWithInvalidToken(ctx context.Context) (context.Context, error) {
	reqCtx := tokenCtx{
		token: "invalid",
	}

	return context.WithValue(ctx, tokenCtxKey{}, reqCtx), nil
}
