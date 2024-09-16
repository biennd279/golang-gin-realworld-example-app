package tests

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)
import "context"
import "github.com/cucumber/godog"

func TestArticle_security(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeArticle_securityScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/article_security.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeArticle_securityScenario(ctx *godog.ScenarioContext) {
	SetupDefaultApplicationScenario(ctx)
	ctx.Step(`^I send a request to "([^"]*)" the article api$`, iSendARequestToActionTheArticleApi)
	ctx.Step(`^the response status code should be a (\d+)$`, theresponsestatuscodeshouldbeastatusCode)
}

func iSendARequestToActionTheArticleApi(ctx context.Context, action string) (context.Context, error) {
	var method string
	var endpoint string

	switch action {
	case "get":
		method = "GET"
		endpoint = "/api/articles/1"
	case "create":
		method = "POST"
		endpoint = "/api/articles/"
	case "update":
		method = "PUT"
		endpoint = "/api/articles/1"
	case "delete":
		method = "DELETE"
		endpoint = "/api/articles/1"
	default:
		return ctx, godog.ErrPending
	}

	req := NewJSONRequest(method, endpoint, nil)

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
func theresponsestatuscodeshouldbeastatusCode(ctx context.Context, statusCode int) (context.Context, error) {
	rsp := ctx.Value(rspCtxKey{}).(response)

	if rsp.statusCode != statusCode {
		return ctx, errors.New(fmt.Sprintf("expected status code %d, got %d", statusCode, rsp.statusCode))
	}
	return ctx, nil
}
