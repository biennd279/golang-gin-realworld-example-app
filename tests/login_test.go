package login_test

import (
	"github.com/cucumber/godog"
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

func iHaveALoginApi() error {
	return godog.ErrPending
}

func iSendARequestToTheLoginApiWithTheFollowingData() error {
	return godog.ErrPending
}

func iShouldGetAResponseFromTheLoginApi() error {
	return godog.ErrPending
}

func theResponseShouldBeAStatusCode(arg1 int) error {
	return godog.ErrPending
}

func theResponseShouldContainAToken() error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have a login api$`, iHaveALoginApi)
	ctx.Step(`^I send a request to the login api with the following data$`, iSendARequestToTheLoginApiWithTheFollowingData)
	ctx.Step(`^I should get a response from the login api$`, iShouldGetAResponseFromTheLoginApi)
	ctx.Step(`^the response should be a (\d+) status code$`, theResponseShouldBeAStatusCode)
	ctx.Step(`^the response should contain a token$`, theResponseShouldContainAToken)
}
