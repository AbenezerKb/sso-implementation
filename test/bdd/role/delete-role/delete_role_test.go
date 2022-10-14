package delete_role

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
)

type deleteRoleTest struct {
	test.TestInstance
	admin, user db.User
	apiTest     src.ApiTest
	role        dto.Role
}

func TestDeleteRole(t *testing.T) {
	r := deleteRoleTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "delete role test", "features/delete_role.feature", r.InitializeScenario)
}

func (r *deleteRoleTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
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

func (r *deleteRoleTest) thereIsARoleWithTheFollowingDetails(roleTable *godog.Table) error {
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

func (r *deleteRoleTest) theFollowingUserHasTheRoleAssigned(userTable *godog.Table) error {
	userJSON, err := r.apiTest.ReadRow(userTable, nil, false)
	if err != nil {
		return err
	}

	var user dto.User
	err = r.apiTest.UnmarshalJSON([]byte(userJSON), &user)
	if err != nil {
		return err
	}

	r.user, err = r.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email: sql.NullString{
			String: user.Email,
			Valid:  true,
		},
		Phone:    user.Phone,
		Password: user.Password,
	})
	if err != nil {
		return err
	}
	return r.PersistDB.AssignRoleForUser(context.Background(), r.user.ID, r.role.Name)
}

func (r *deleteRoleTest) iRequestToDeleteTheRole(role string) error {
	r.apiTest.URL = r.apiTest.URL + role
	r.apiTest.SendRequest()
	return nil
}

func (r *deleteRoleTest) theRoleShouldBeDeleted() error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	_, err := r.DB.GetRoleByName(context.Background(), r.role.Name)
	if err == nil {
		return fmt.Errorf("expected to not find the deleted role")
	}

	rows, err := r.Conn.Query(context.Background(), "SELECT * FROM casbin_rule WHERE v0 = $1", r.role.Name)
	if err == nil {
		return err
	}
	if rows.Next() {
		return fmt.Errorf("expected to not find permissions associated with the deleted role")
	}
	return nil
}

func (r *deleteRoleTest) theUserShouldNoLongerHaveThatRoleAssigned() error {
	rows, err := r.Conn.Query(context.Background(), "SELECT * FROM casbin_rule WHERE v1 = $1", r.role.Name)
	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf("expected to not find any users associated with the deleted role")
	}
	return nil
}

func (r *deleteRoleTest) myRequestShouldFailWith(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (r *deleteRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/roles/"
		r.apiTest.Method = http.MethodDelete

		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteRole(ctx, r.role.Name)
		_, _ = r.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v0 = $1", r.role.Name)
		_, _ = r.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v1 = $1", r.role.Name)
		_, _ = r.DB.DeleteUser(ctx, r.admin.ID)
		_, _ = r.DB.DeleteUser(ctx, r.user.ID)
		_ = r.GrantRoleAfterFunc()
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, r.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I request to delete the role "([^"]*)"$`, r.iRequestToDeleteTheRole)
	ctx.Step(`^my request should fail with "([^"]*)"$`, r.myRequestShouldFailWith)
	ctx.Step(`^the following user has the role assigned$`, r.theFollowingUserHasTheRoleAssigned)
	ctx.Step(`^the role should be deleted$`, r.theRoleShouldBeDeleted)
	ctx.Step(`^the user should no longer have that role assigned$`, r.theUserShouldNoLongerHaveThatRoleAssigned)
	ctx.Step(`^there is a role with the following details:$`, r.thereIsARoleWithTheFollowingDetails)
}
