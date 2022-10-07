package get_role

import (
	"context"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils/collection"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getRoleTest struct {
	test.TestInstance
	apiTest src.ApiTest
	roles   []dto.Role
	admin   db.User
}

func TestGetRole(t *testing.T) {
	g := &getRoleTest{}
	g.apiTest.Method = http.MethodGet
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get role test", "features/get_role.feature", g.InitializeScenario)
}

func (g *getRoleTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	g.admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, g.GrantRoleAfterFunc, err = g.GrantRoleForUserWithAfter(g.admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}

	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)

	return nil
}
func (g *getRoleTest) theFollowingRolesAreRegisteredOnTheSystem(roles *godog.Table) error {
	rolesJSON, err := g.apiTest.ReadRows(roles, []src.Type{
		{
			Column: "permissions",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}

	var rolesData []dto.Role
	err = g.apiTest.UnmarshalJSON([]byte(rolesJSON), &rolesData)
	if err != nil {
		return err
	}

	for _, v := range rolesData {
		role, err := g.PersistDB.CreateRoleTX(context.Background(), v.Name, v.Permissions)
		if err != nil {
			return err
		}
		g.roles = append(g.roles, role)
	}

	return nil
}

func (g *getRoleTest) iRequestToGetARoleBy(role string) error {
	g.apiTest.URL += role
	g.apiTest.SendRequest()
	return nil
}

func (g *getRoleTest) iShouldGetTheFollowingRole(role *godog.Table) error {
	var responseRole dto.Role

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	roleJSON, err := g.apiTest.ReadRow(role, []src.Type{
		{
			Column: "permissions",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}

	var roleData dto.Role
	err = g.apiTest.UnmarshalJSON([]byte(roleJSON), &roleData)
	if err != nil {
		return err
	}

	err = g.apiTest.UnmarshalResponseBodyPath("data", &responseRole)
	if err != nil {
		return err
	}

	for _, v := range responseRole.Permissions {
		if !collection.Contains(v, roleData.Permissions) {
			return fmt.Errorf("expected to get permission %s in %v", v, roleData.Permissions)
		}
	}

	return nil
}

func (g *getRoleTest) myRequestShouldFailWith(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return g.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (g *getRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	g.apiTest.URL = "/v1/roles/"
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.admin.ID)

		for _, v := range g.roles {
			_, _ = g.DB.DeleteRole(ctx, v.Name)
			_, _ = g.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v0 = $1", v.Name)
		}
		_ = g.GrantRoleAfterFunc()

		return ctx, nil
	})
	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get a role by "([^"]*)"$`, g.iRequestToGetARoleBy)
	ctx.Step(`^I should get the following role$`, g.iShouldGetTheFollowingRole)
	ctx.Step(`^my request should fail with "([^"]*)"$`, g.myRequestShouldFailWith)
	ctx.Step(`^The following roles are registered on the system$`, g.theFollowingRolesAreRegisteredOnTheSystem)
}
