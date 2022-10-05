package delete_scope

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type deleteScopeTest struct {
	test.TestInstance
	apiTest   src.ApiTest
	Admin     db.User
	scopes    []db.Scope
	scopeName string
}

func TestDeleteScope(t *testing.T) {
	d := deleteScopeTest{}
	d.TestInstance = test.Initiate("../../../../")
	d.apiTest.InitializeTest(t, "delete scope", "features/delete_scope.feature", d.InitializeScenario)
}

func (d *deleteScopeTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
	body, err := d.apiTest.ReadRow(adminCredential, nil, false)
	if err != nil {
		return err
	}

	adminValue := dto.User{}
	err = d.apiTest.UnmarshalJSON([]byte(body), &adminValue)
	if err != nil {
		return err
	}

	d.Admin, err = d.AuthenticateWithParam(adminValue)
	if err != nil {
		return err
	}
	d.apiTest.SetHeader("Authorization", "Bearer "+d.AccessToken)
	return d.GrantRoleForUser(d.Admin.ID.String(), adminCredential)
}

func (d *deleteScopeTest) iDeleteTheScopeWithName(scopeName string) error {
	d.apiTest.URL += scopeName
	d.scopeName = scopeName
	d.apiTest.SendRequest()
	return nil
}

func (d *deleteScopeTest) theScopeShouldBeDeleted() error {
	if err := d.apiTest.AssertStatusCode(http.StatusNoContent); err != nil {
		return err
	}

	_, err := d.DB.GetScope(context.Background(), d.scopeName)
	if err == nil {
		return fmt.Errorf("scope is not deleted")
	}

	if !sqlcerr.Is(err, sqlcerr.ErrNoRows) {
		return err
	}

	return nil
}

func (d *deleteScopeTest) theDeleteShouldFailWithMessage(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return d.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (d *deleteScopeTest) thereAreScopesWithTheFollowingDetails(scopes *godog.Table) error {
	scopesData, err := d.apiTest.ReadRows(scopes, nil, false)
	if err != nil {
		return err
	}
	var scopesStruct []dto.Scope
	if err := d.apiTest.UnmarshalJSONAt([]byte(scopesData), "", &scopesStruct); err != nil {
		return err
	}
	for _, scope := range scopesStruct {
		savedScope, err := d.DB.CreateScope(context.Background(), db.CreateScopeParams{
			Name:        scope.Name,
			Description: scope.Description,
			ResourceServerName: sql.NullString{
				String: scope.ResourceServerName,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
		d.scopes = append(d.scopes, savedScope)
	}
	return nil
}

func (d *deleteScopeTest) InitializeScenario(ctx *godog.ScenarioContext) {
	d.apiTest.URL = "/v1/oauth/scopes/"
	d.apiTest.Method = http.MethodDelete
	d.apiTest.SetHeader("Content-Type", "application/json")
	d.apiTest.InitializeServer(d.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.DeleteUser(ctx, d.Admin.ID)

		for _, scope := range d.scopes {
			_, _ = d.DB.DeleteScope(ctx, scope.Name)
		}

		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, d.iAmLoggedInAsAdminUser)
	ctx.Step(`^I delete the scope with name "([^"]*)"$`, d.iDeleteTheScopeWithName)
	ctx.Step(`^The delete should fail with message "([^"]*)"$`, d.theDeleteShouldFailWithMessage)
	ctx.Step(`^The scope should be deleted$`, d.theScopeShouldBeDeleted)
	ctx.Step(`^There are scope\'s with the following details$`, d.thereAreScopesWithTheFollowingDetails)
}
