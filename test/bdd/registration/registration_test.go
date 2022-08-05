package registration

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

type registrationTest struct {
	apiTest src.ApiTest
	server  *gin.Engine
	db      *db.Queries
}

func TestRegistertion(t *testing.T) {

	a := &registrationTest{}
	a.server, a.db = test.GetServer("../../../")

	a.apiTest.InitializeTest(t, "Login test", "features/registration.feature", a.InitializeScenario)
}

func (r *registrationTest) iFillTheFormWithTheFollowingDetails(userForm *godog.Table) error {
	body, err := r.apiTest.ReadRow(userForm, nil, false)
	if err != nil {
		return err
	}
	r.apiTest.Body = body
	return nil
}

func (r *registrationTest) iSubmitTheRegistrationForm() error {
	r.apiTest.SendRequest()
	return nil
}

func (r *registrationTest) iWillHaveANewAccount() error {
	if err := r.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}

	return nil
}

func (r *registrationTest) theRegistrationShouldFailWith(msg string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if err := r.apiTest.AssertPathValue(fmt.Sprintf(`{"message":"%s"}`, msg), "message", string(r.apiTest.ResponseBody), "error.field_error.0.description"); err != nil {
		return err
	}

	return nil
}

func (r *registrationTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {

		r.apiTest.URL = "/v1/register"
		r.apiTest.Method = http.MethodPost
		r.apiTest.Headers = map[string]string{}
		r.apiTest.Headers["Content-Type"] = "application/json"
		r.apiTest.InitializeServer(r.server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})

	ctx.Step(`^I fill the form with the following details$`, r.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I submit the registration form$`, r.iSubmitTheRegistrationForm)
	ctx.Step(`^I will have a new account$`, r.iWillHaveANewAccount)
	ctx.Step(`^the registration should fail with "([^"]*)"$`, r.theRegistrationShouldFailWith)
}
