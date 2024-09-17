package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"net/http/httptest"
	"testing"
)

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

func iSendARequestToTheLoginApiWithValidCredentials(ctx context.Context) (context.Context, error) {
	loginReq := CreateLoginRequest("test@gmail.com", "password")

	w := httptest.NewRecorder()

	app := ctx.Value(appCtxKey{}).(*appContext)

	app.r.ServeHTTP(w, loginReq)

	rsp := responseCtx{
		statusCode: w.Code,
		body:       w.Body.String(),
		headers:    w.Header(),
	}

	return context.WithValue(ctx, rspCtxKey{}, rsp), nil
}

func theResponseShouldBeAStatusCode(ctx context.Context, arg1 int) (context.Context, error) {
	rsp := ctx.Value(rspCtxKey{}).(responseCtx)

	if rsp.statusCode != arg1 {
		return ctx, errors.New(fmt.Sprintf("expected status code %d, got %d", arg1, rsp.statusCode))
	}

	return ctx, nil
}

func theResponseShouldContainAToken(ctx context.Context) error {
	rsp := ctx.Value(rspCtxKey{}).(responseCtx)
	var loginResponse userLoginResponse
	err := json.Unmarshal([]byte(rsp.body), &loginResponse)
	if err != nil {
		return err
	}

	if loginResponse.User.Token == "" {
		return errors.New("expected token in responseCtx")
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	SetupDefaultApplicationScenario(ctx)
	ctx.Step(`^the response should be a (\d+) status code$`, theResponseShouldBeAStatusCode)
	ctx.Step(`^the response should contain a token$`, theResponseShouldContainAToken)
	ctx.Step(`^I send a request to the login api with valid credentials$`, iSendARequestToTheLoginApiWithValidCredentials)
}
