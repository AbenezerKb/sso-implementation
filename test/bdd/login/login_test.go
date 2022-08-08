package login

import (
	"context"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type loginTest struct {
	apiTest src.ApiTest
	server  *gin.Engine
	db      *db.Queries
}

func TestLogin(t *testing.T) {

	a := &loginTest{}
	a.server, a.db = test.GetServer("../../../config")

	a.apiTest.InitializeTest(t, "Login test", "features/login.feature", a.InitializeScenario)
}

func (l *loginTest) iFillTheFollowingDetails(loginInfo *godog.Table) error {
	body, err := l.apiTest.ReadRow(loginInfo, nil, false)
	if err != nil {
		return err
	}
	l.apiTest.Body = body
	return nil
}

func (l *loginTest) iSubmitTheRegistrationForm() error {
	l.apiTest.SendRequest()
	return nil
}

func (l *loginTest) iWillBeLoggedInSecurelyToMyAccount() error {
	if err := l.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}

	return nil
}

func (l *loginTest) theLoginShouldFailWith(msg string) error {
	if err := l.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if err := l.apiTest.AssertPathValue(fmt.Sprintf(`{"message":"%s"}`, msg), "message", string(l.apiTest.ResponseBody), "error.field_error.0.description"); err != nil {
		return err
	}
	return nil
}

func (l *loginTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {

		l.apiTest.URL = "/v1/login"
		l.apiTest.Method = http.MethodPost
		l.apiTest.Headers = map[string]string{}
		l.apiTest.Headers["Content-Type"] = "application/json"
		l.apiTest.InitializeServer(l.server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})
	ctx.Step(`^I fill the following details$`, l.iFillTheFollowingDetails)
	ctx.Step(`^I submit the registration form$`, l.iSubmitTheRegistrationForm)
	ctx.Step(`^I will be logged in securely to my account$`, l.iWillBeLoggedInSecurelyToMyAccount)
	ctx.Step(`^the login should fail with "([^"]*)"$`, l.theLoginShouldFailWith)
}
