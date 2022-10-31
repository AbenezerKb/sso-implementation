package get_user_permissions

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/platform/utils/collection"
	"sso/test"
	"testing"
)

type getUserPermissionsTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	user        db.User
	permissions []string
}

func TestPermissions(t *testing.T) {
	g := getUserPermissionsTest{}

	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.URL = "/v1/profile/permissions"
	g.apiTest.Method = http.MethodGet
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get user permissions", "features/get_user_permissions.feature", g.InitializeScenario)
}

func (g *getUserPermissionsTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error

	g.user, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}

	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)

	permissions, err := g.apiTest.ReadCell(adminCredentials, "role", &src.Type{Kind: src.Array})
	if err != nil {
		return err
	}

	permissionNames, ok := permissions.([]string)
	if !ok {
		return fmt.Errorf("couldn't scan permissions from table")
	}

	g.permissions = permissionNames
	_, g.GrantRoleAfterFunc, err = g.GrantRoleForUserWithAfter(g.user.ID.String(), adminCredentials)
	return err
}

func (g *getUserPermissionsTest) iRequestToGetMyPermissions() error {
	g.apiTest.SendRequest()
	return nil
}

func (g *getUserPermissionsTest) iShouldGetAllMyPermissions() error {
	var responsePermissions []string

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	err := g.apiTest.UnmarshalResponseBodyPath("data", &responsePermissions)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(len(responsePermissions), len(g.permissions)); err != nil {
		return err
	}

	for _, v := range g.permissions {
		if !collection.Contains(v, responsePermissions) {
			return fmt.Errorf("expected permisson: %s", v)
		}
	}

	return nil
}

func (g *getUserPermissionsTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.user.ID)
		_ = g.DB.RemoveInternalRefreshTokenByUserID(ctx, g.user.ID)
		_ = g.GrantRoleAfterFunc()
		return ctx, nil
	})

	ctx.Step(`^I am logged in with the following credentials$`, g.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I request to get my permissions$`, g.iRequestToGetMyPermissions)
	ctx.Step(`^I should get all my permissions$`, g.iShouldGetAllMyPermissions)
}
