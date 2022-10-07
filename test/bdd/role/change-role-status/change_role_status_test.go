package change_role_status

import (
	"context"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type changeRoleStatusTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Admin   db.User
	role    dto.Role
}

func TestUpdateClientStatus(t *testing.T) {
	r := changeRoleStatusTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.InitializeTest(t, "change role status", "features/change_role_status.feature", r.InitializeScenario)
}

func (r *changeRoleStatusTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	r.Admin, err = r.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, r.GrantRoleAfterFunc, err = r.GrantRoleForUserWithAfter(r.Admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}
	return nil
}

func (r *changeRoleStatusTest) thereIsARoleWithTheFollowingDetails(roleTable *godog.Table) error {

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

	r.apiTest.URL += r.role.Name + "/status"

	return nil
}

func (r *changeRoleStatusTest) iUpdateTheRolesStatusTo(updatedStatus string) error {
	r.apiTest.SetBodyMap(map[string]interface{}{
		"status": updatedStatus,
	})
	r.apiTest.SendRequest()
	return nil
}

func (r *changeRoleStatusTest) theRoleStatusShouldUpdateTo(updatedStatus string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	updatedRoleData, err := r.DB.GetRoleByName(context.Background(), r.role.Name)
	if err != nil {
		return err
	}

	if err := r.apiTest.AssertEqual(updatedRoleData.Status, updatedStatus); err != nil {
		return err
	}
	return nil
}

func (r *changeRoleStatusTest) thenIShouldGetRoleNotFoundErrorWithMessage(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return r.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (r *changeRoleStatusTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (r *changeRoleStatusTest) thereIsRoleWithName(roleName string) error {
	r.apiTest.URL += roleName + "/status"
	return nil
}

func (r *changeRoleStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
	r.apiTest.URL = "/v1/roles/"
	r.apiTest.Method = http.MethodPatch
	r.apiTest.InitializeServer(r.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteRole(ctx, r.role.Name)
		_, _ = r.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v0 = $1", r.role.Name)
		_, _ = r.DB.DeleteUser(ctx, r.Admin.ID)
		_ = r.GrantRoleAfterFunc()
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, r.iAmLoggedInAsAdminUser)
	ctx.Step(`^I update the role\'s status to "([^"]*)"$`, r.iUpdateTheRolesStatusTo)
	ctx.Step(`^the role status should update to "([^"]*)"$`, r.theRoleStatusShouldUpdateTo)
	ctx.Step(`^Then I should get role not found error with message "([^"]*)"$`, r.thenIShouldGetRoleNotFoundErrorWithMessage)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, r.thenIShouldGetErrorWithMessage)
	ctx.Step(`^there is role with name "([^"]*)"$`, r.thereIsRoleWithName)
	ctx.Step(`^there is a role with the following details:$`, r.thereIsARoleWithTheFollowingDetails)
}
