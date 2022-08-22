package create

import (
	"context"
	"net/http"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type scopeCreateTest struct {
	test.TestInstance
	apiTest src.ApiTest
	scope   dto.Scope
}

func TestScopeCreation(t *testing.T) {
	s := &scopeCreateTest{}
	s.TestInstance = test.Initiate("../../../../")
	s.apiTest.InitializeTest(t, "scope creation test", "features/scope.feature", s.InitializeScenario)
}

func (s *scopeCreateTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	err := s.Authenicate(adminCredentials)
	if err != nil {
		return err
	}

	return nil
}

func (s *scopeCreateTest) iCreateTheScope() error {
	s.apiTest.SetHeader("Authorization", "Bearer "+s.AccessToken)
	s.apiTest.SendRequest()
	return nil
}

func (s *scopeCreateTest) iFillTheFormWithFollowingFields(form *godog.Table) error {
	scope, err := s.apiTest.ReadRow(form, nil, false)
	if err != nil {
		return err
	}
	s.apiTest.Body = scope
	return nil
}

func (s *scopeCreateTest) iShouldHaveNewScope() error {
	if err := s.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	return s.apiTest.UnmarshalResponseBodyPath("data", &s.scope)
}

func (s *scopeCreateTest) theCreationShouldFailWith(message string) error {
	if err := s.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return s.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)

}

func (s *scopeCreateTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.apiTest.URL = "/v1/oauth/scope"
		s.apiTest.Method = http.MethodPost
		s.apiTest.SetHeader("Content-Type", "application/json")
		s.apiTest.InitializeServer(s.Server)
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// // delete scope
		_, _ = s.DB.DeleteScope(ctx, s.scope.Name)
		_, err = s.DB.DeleteUser(ctx, s.User.ID)
		return ctx, err
	})
	ctx.Step(`^I am logged in as admin user$`, s.iAmLoggedInAsAdminUser)
	ctx.Step(`^I create the scope$`, s.iCreateTheScope)
	ctx.Step(`^I fill the form with following fields:$`, s.iFillTheFormWithFollowingFields)
	ctx.Step(`^I should have new scope$`, s.iShouldHaveNewScope)
	ctx.Step(`^The creation should fail with "([^"]*)"$`, s.theCreationShouldFailWith)
}
