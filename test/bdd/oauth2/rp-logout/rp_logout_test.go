package logout

import (
	"context"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type rpLogoutTest struct {
	test.TestInstance
	apiTest src.ApiTest
	User    db.User
}

func TestLogout(t *testing.T) {
	r := &rpLogoutTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.InitializeTest(t, "Logout test", "features/rp_logout.feature", r.InitializeScenario)
}

func (r *rpLogoutTest) iAmRegisteredOnTheSystem() error {
	return godog.ErrPending
}

func (r *rpLogoutTest) iHaveId_token() error {
	return godog.ErrPending
}

func (r *rpLogoutTest) iHaveTheFollowingDetails(arg1 *godog.Table) error {
	return godog.ErrPending
}

func (r *rpLogoutTest) iHaveTheFollowingInvalid_requestDetails(arg1 *godog.Table) error {
	return godog.ErrPending
}

func (r *rpLogoutTest) iRequestToLogout() error {
	return godog.ErrPending
}

func (r *rpLogoutTest) iShouldBeRedirectedToWithTheFollowingQueryParams(arg1 string, arg2 *godog.Table) error {
	return godog.ErrPending
}

func (r *rpLogoutTest) theUserIsRegisteredOnTheSystem() error {
	return godog.ErrPending
}

func (r *rpLogoutTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/oauth/logout"
		r.apiTest.Method = http.MethodGet
		r.apiTest.SetHeader("Content-Type", "application/json")
		r.apiTest.InitializeServer(r.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(ctx, r.User.ID)
		return ctx, err
	})

	ctx.Step(`^I am  registered on the system$`, r.iAmRegisteredOnTheSystem)
	ctx.Step(`^I have id_token$`, r.iHaveId_token)
	ctx.Step(`^I have the following details:$`, r.iHaveTheFollowingDetails)
	ctx.Step(`^I have the following invalid_request details:$`, r.iHaveTheFollowingInvalid_requestDetails)
	ctx.Step(`^I request to logout$`, r.iRequestToLogout)
	ctx.Step(`^I should be redirected to "([^"]*)" with the following query params:$`, r.iShouldBeRedirectedToWithTheFollowingQueryParams)
	ctx.Step(`^the user is registered on the system$`, r.theUserIsRegisteredOnTheSystem)
}
