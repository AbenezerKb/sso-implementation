package assign_role

import (
	"context"
	"database/sql"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
)

type assignRoleTest struct {
	test.TestInstance
	apiTest         src.ApiTest
	admin, user     db.User
	adminRole, role *dto.Role
}

func TestAssignRole(t *testing.T) {
	a := assignRoleTest{}
	a.apiTest.URL = "/v1/users" // initial prefix
	a.apiTest.Method = http.MethodPatch
	a.TestInstance = test.Initiate("../../../../")
	a.apiTest.InitializeServer(a.Server)
	a.apiTest.InitializeTest(t, "assign role test", "features/assign_role.feature", a.InitializeScenario)
}

// backgrounds
func (a *assignRoleTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error
	a.admin, err = a.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	a.adminRole, a.GrantRoleAfterFunc, err = a.GrantRoleForUserWithAfter(a.admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}
	return nil
}

func (a *assignRoleTest) theFollowingRoleIsRegisteredOnTheSystem(roleTable *godog.Table) error {
	roleJSON, err := a.apiTest.ReadRow(roleTable, []src.Type{
		{
			Column: "permissions",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}

	var role dto.Role
	err = a.apiTest.UnmarshalJSON([]byte(roleJSON), &role)
	if err != nil {
		return err
	}

	roleDB, err := a.PersistDB.CreateRoleTX(context.Background(), role.Name, role.Permissions)
	if err != nil {
		return err
	}
	a.role = &roleDB
	return nil
}

func (a *assignRoleTest) theFollowingUserIsRegisteredOnTheSystem(userTable *godog.Table) error {
	userJSON, err := a.apiTest.ReadRow(userTable, nil, false)
	if err != nil {
		return err
	}

	var user dto.User
	err = a.apiTest.UnmarshalJSON([]byte(userJSON), &user)
	if err != nil {
		return err
	}

	a.user, err = a.DB.CreateUser(context.Background(), db.CreateUserParams{
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
	return err
}

// when
func (a *assignRoleTest) iRequestToAssignAsRoleForTheUser(role string) error {
	a.apiTest.SetHeader("Authorization", "Bearer "+a.AccessToken)
	a.apiTest.SetBodyValue("role", role)
	a.apiTest.URL = a.apiTest.URL + "/" + a.user.ID.String() + "/role"
	a.apiTest.SendRequest()
	return nil
}

// then
func (a *assignRoleTest) theRoleShouldBeAssignedToTheUser() error {
	if err := a.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	var roleName string
	row := a.Conn.QueryRow(context.Background(), "SELECT v1 FROM casbin_rule WHERE v0 = $1", a.user.ID)
	if err := row.Scan(&roleName); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(roleName, a.role.Name); err != nil {
		return err
	}
	return nil
}

func (a *assignRoleTest) myRequestShouldFailWithAnd(message, fieldError string) error {
	if err := a.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if message != "" {
		if err := a.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
			return err
		}
	}
	if fieldError != "" {
		if err := a.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", fieldError); err != nil {
			return err
		}
	}

	return nil
}

func (a *assignRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = a.DB.DeleteUser(ctx, a.user.ID)
		_, _ = a.DB.DeleteUser(ctx, a.admin.ID)
		_, _ = a.DB.DeleteRole(ctx, a.role.Name)
		_, _ = a.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v0 = $1", a.role.Name)
		_ = a.GrantRoleAfterFunc()
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, a.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I request to assign "([^"]*)" as role for the user$`, a.iRequestToAssignAsRoleForTheUser)
	ctx.Step(`^my request should fail with "([^"]*)" and "([^"]*)"$`, a.myRequestShouldFailWithAnd)
	ctx.Step(`^The following role is registered on the system$`, a.theFollowingRoleIsRegisteredOnTheSystem)
	ctx.Step(`^The following user is registered on the system$`, a.theFollowingUserIsRegisteredOnTheSystem)
	ctx.Step(`^the role should be assigned to the user$`, a.theRoleShouldBeAssignedToTheUser)
}
