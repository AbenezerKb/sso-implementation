package update_scope

import (
	"context"
	"database/sql"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateScopeTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	Admin       db.User
	scope       db.Scope
	updateScope dto.UpdateScopeParam
}

func TestUpdateScope(t *testing.T) {
	u := updateScopeTest{}
	u.TestInstance = test.Initiate("../../../../")

	u.apiTest.Method = http.MethodPut
	u.apiTest.InitializeServer(u.Server)

	u.apiTest.InitializeTest(t, "update scope", "features/update_scope.feature", u.InitializeScenario)
}

func (u *updateScopeTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
	body, err := u.apiTest.ReadRow(adminCredential, nil, false)
	if err != nil {
		return err
	}

	adminValue := dto.User{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &adminValue)
	if err != nil {
		return err
	}

	u.Admin, err = u.AuthenticateWithParam(adminValue)
	if err != nil {
		return err
	}
	u.apiTest.SetHeader("Authorization", "Bearer "+u.AccessToken)
	return u.GrantRoleForUser(u.Admin.ID.String(), adminCredential)
}

func (u *updateScopeTest) thereIsScopeWithTheFollowingDetails(scope *godog.Table) error {
	scopeData, err := u.apiTest.ReadRow(scope, nil, false)
	if err != nil {
		return err
	}
	var scopeStruct dto.Scope
	if err := u.apiTest.UnmarshalJSONAt([]byte(scopeData), "", &scopeStruct); err != nil {
		return err
	}
	u.scope, err = u.DB.CreateScope(context.Background(), db.CreateScopeParams{
		Name:        scopeStruct.Name,
		Description: scopeStruct.Description,
		ResourceServerName: sql.NullString{
			String: scopeStruct.ResourceServerName,
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}

	u.apiTest.URL += u.scope.Name

	return nil
}

func (u *updateScopeTest) thereIsScope(name string) error {
	u.apiTest.URL += name
	return nil
}

func (u *updateScopeTest) iFillTheFollowingDetails(scopeUpdate *godog.Table) error {
	body, err := u.apiTest.ReadRow(scopeUpdate, nil, false)
	if err != nil {
		return err
	}

	updateScopeValues := dto.UpdateScopeParam{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &updateScopeValues)
	if err != nil {
		return err
	}

	u.updateScope = updateScopeValues

	u.apiTest.Body = body
	return nil
}

func (u *updateScopeTest) iUpdateScope() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *updateScopeTest) theScopeShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	fetchedScope, err := u.DB.GetScope(context.Background(), u.scope.Name)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedScope.Description, u.updateScope.Description); err != nil {
		return err
	}

	return nil
}

func (u *updateScopeTest) theUpdateShouldFailWithFieldErrorDescription(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (u *updateScopeTest) theUpdateShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (u *updateScopeTest) InitializeScenario(ctx *godog.ScenarioContext) {
	u.apiTest.URL = "/v1/oauth/scopes/"
	u.apiTest.SetHeader("Content-Type", "application/json")

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteScope(ctx, u.scope.Name)
		_, _ = u.DB.DeleteUser(ctx, u.Admin.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, u.iAmLoggedInAsAdminUser)
	ctx.Step(`^I fill the following details$`, u.iFillTheFollowingDetails)
	ctx.Step(`^I update scope$`, u.iUpdateScope)
	ctx.Step(`^The scope should be updated$`, u.theScopeShouldBeUpdated)
	ctx.Step(`^The update should fail with field error description "([^"]*)"$`, u.theUpdateShouldFailWithFieldErrorDescription)
	ctx.Step(`^The update should fail with message "([^"]*)"$`, u.theUpdateShouldFailWithMessage)
	ctx.Step(`^there is scope "([^"]*)"$`, u.thereIsScope)
	ctx.Step(`^there is scope with the following details$`, u.thereIsScopeWithTheFollowingDetails)
}
