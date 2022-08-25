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

type logoutTest struct {
	test.TestInstance
	apiTest src.ApiTest
	User    db.User
}

func TestLogout(t *testing.T) {
	l := &logoutTest{}
	l.TestInstance = test.Initiate("../../../../")
	l.apiTest.InitializeTest(t, "Logout test", "features/logout.feature", l.InitializeScenario)
}

func (l *logoutTest) iAmALogedinUserWithTheFollowingDetails(userCredentials *godog.Table) error {
	var err error
	l.User, err = l.Authenticate(userCredentials)
	if err != nil {
		return err
	}

	return nil
}

func (l *logoutTest) iLogout() error {
	l.apiTest.SetHeader("Authorization", "Bearer "+l.AccessToken)
	l.apiTest.SendRequest()
	return nil
}

func (l *logoutTest) iShouldSuccessfullyLogoutOfTheSystem() error {
	if err := l.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	return nil
}

func (l *logoutTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		l.apiTest.URL = "/v1/logout"
		l.apiTest.Method = http.MethodGet
		l.apiTest.SetHeader("Content-Type", "application/json")
		l.apiTest.InitializeServer(l.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = l.DB.DeleteUser(ctx, l.User.ID)
		return ctx, err
	})

	ctx.Step(`^I am a logedin  user with the following details:$`, l.iAmALogedinUserWithTheFollowingDetails)
	ctx.Step(`^I should Successfully logout of the system$`, l.iShouldSuccessfullyLogoutOfTheSystem)
	ctx.Step(`^I logout$`, l.iLogout)
}
