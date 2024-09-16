package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"github.com/gothinkster/golang-gin-realworld-example-app/users"
	"golang.org/x/crypto/bcrypt"
	"net/http/httptest"
	"testing"
)

type userLoginResponse struct {
	User struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Bio      string  `json:"bio"`
		Image    *string `json:"image"`
		Token    string  `json:"token"`
	} `json:"user"`
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/login_security.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

//func sendRequestToTheLoginApi(r *gin.Engine) (*http.Response, error) {
//	return nil, godog.ErrPending
//
//}

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

func iSendARequestToTheLoginApiWithValidCredentials(ctx context.Context) (context.Context, error) {

	param := gin.H{
		"user": gin.H{
			"email":    "test@gmail.com",
			"password": "password",
		},
	}

	req := NewJSONRequest("POST", "/api/users/login", param)

	w := httptest.NewRecorder()

	app := ctx.Value(appCtxKey{}).(*appContext)

	app.r.ServeHTTP(w, req)

	rsp := response{
		statusCode: w.Code,
		body:       w.Body.String(),
		headers:    w.Header(),
	}

	return context.WithValue(ctx, rspCtxKey{}, rsp), nil
}

func theResponseShouldBeAStatusCode(ctx context.Context, arg1 int) (context.Context, error) {
	rsp := ctx.Value(rspCtxKey{}).(response)

	if rsp.statusCode != arg1 {
		return ctx, errors.New(fmt.Sprintf("expected status code %d, got %d", arg1, rsp.statusCode))
	}

	return ctx, nil
}

func theResponseShouldContainAToken(ctx context.Context) error {
	rsp := ctx.Value(rspCtxKey{}).(response)
	var loginResponse userLoginResponse
	err := json.Unmarshal([]byte(rsp.body), &loginResponse)
	if err != nil {
		return err
	}

	if loginResponse.User.Token == "" {
		return errors.New("expected token in response")
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
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

	ctx.Step(`^the response should be a (\d+) status code$`, theResponseShouldBeAStatusCode)
	ctx.Step(`^the response should contain a token$`, theResponseShouldContainAToken)
	ctx.Step(`^I send a request to the login api with valid credentials$`, iSendARequestToTheLoginApiWithValidCredentials)
	ctx.Step(`^I have a invalid username and password is "([^"]*)" and "([^"]*)"$`, iHaveAInvalidUsernameAndPasswordIsInvalidAndInvalid)
	ctx.Step(`^I have a valid email and password is "([^"]*)" and "([^"]*)"$`, iHaveAValidEmailAndPasswordIs)
}
