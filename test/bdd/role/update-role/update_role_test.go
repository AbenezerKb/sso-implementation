package update_role

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils/collection"
	"sso/test"
	"testing"
)

type updateRoleTest struct {
	test.TestInstance
	apiTest           src.ApiTest
	role              dto.Role
	updatePermissions []string
	admin             db.User
}

func TestUpdateRole(t *testing.T) {
	r := updateRoleTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "update role test", "features/update_role.feature", r.InitializeScenario)
}

func (r *updateRoleTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error
	r.admin, err = r.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, r.GrantRoleAfterFunc, err = r.GrantRoleForUserWithAfter(r.admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}
	r.apiTest.SetHeader("Authorization", "Bearer "+r.AccessToken)
	return nil
}

func (r *updateRoleTest) thereIsARoleWithTheFollowingDetails(roleTable *godog.Table) error {
	roleJSON, err := r.apiTest.ReadRow(roleTable, []src.Type{
		{
			Column: "permissions",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}

	var roleData dto.Role
	err = r.apiTest.UnmarshalJSON([]byte(roleJSON), &roleData)
	if err != nil {
		return err
	}

	r.role, err = r.PersistDB.CreateRoleTX(context.Background(), roleData.Name, roleData.Permissions)
	if err != nil {
		return err
	}

	return nil
}

func (r *updateRoleTest) iRequestToUpdateWithTheFollowingPermissions(roleName string, permissionsTable *godog.Table) error {
	r.apiTest.URL = r.apiTest.URL + roleName
	permissions, err := r.apiTest.ReadCell(permissionsTable, "permissions", &src.Type{Kind: src.Array})
	if err != nil {
		return err
	}
	var ok bool
	r.updatePermissions, ok = permissions.([]string)
	if !ok {
		return fmt.Errorf("error while reading permissions from table")
	}
	r.apiTest.SetBodyValue("permissions", r.updatePermissions)
	r.apiTest.SendRequest()
	return nil
}

func (r *updateRoleTest) myRequestShouldFailWithAnd(message, fieldError string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if message != "" {
		if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
			return err
		}
	}
	if fieldError != "" {
		if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", fieldError); err != nil {
			return err
		}
	}

	return nil
}

func (r *updateRoleTest) myRequestShouldFailWithNoRoleFound(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (r *updateRoleTest) theRoleShouldBeUpdated() error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	var responsePermission dto.Role
	err := r.apiTest.UnmarshalResponseBodyPath("data", &responsePermission)
	if err != nil {
		return err
	}
	if err := r.apiTest.AssertEqual(len(responsePermission.Permissions), len(r.updatePermissions)); err != nil {
		return err
	}
	for _, v := range responsePermission.Permissions {
		if !collection.Contains(v, r.updatePermissions) {
			return fmt.Errorf("expected to get permission %s in %v", v, r.updatePermissions)
		}
	}
	return nil
}
func (r *updateRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/roles/"
		r.apiTest.Method = http.MethodPut

		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteRole(ctx, r.role.Name)
		_, _ = r.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v0 = $1", r.role.Name)
		_, _ = r.DB.DeleteUser(ctx, r.admin.ID)
		_ = r.GrantRoleAfterFunc()
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, r.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I request to update "([^"]*)" with the following permissions$`, r.iRequestToUpdateWithTheFollowingPermissions)
	ctx.Step(`^the role should be updated$`, r.theRoleShouldBeUpdated)
	ctx.Step(`^my request should fail with "([^"]*)" and "([^"]*)"$`, r.myRequestShouldFailWithAnd)
	ctx.Step(`^there is a role with the following details$`, r.thereIsARoleWithTheFollowingDetails)
	ctx.Step(`^my request should fail with no role found "([^"]*)"$`, r.myRequestShouldFailWithNoRoleFound)
}
