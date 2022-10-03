package get_permission_list

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/permissions"
	"sso/platform/utils/collection"
	"sso/test"
	"testing"
)

type getPermissionsTest struct {
	test.TestInstance
	apiTest  src.ApiTest
	Admin    db.User
	category string
}

func TestGetPermissionList(t *testing.T) {
	g := getPermissionsTest{}
	g.apiTest.URL = "/v1/roles/permissions"
	g.apiTest.Method = http.MethodGet
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get permissions test", "features/get_permission_list.feature", g.InitializeScenario)
}

func (g *getPermissionsTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error
	g.Admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	return g.GrantRoleForUser(g.Admin.ID.String(), adminCredentials)
}

func (g *getPermissionsTest) iRequestToGetAllPermissionsWithCategory(category string) error {
	g.category = category
	g.apiTest.SetQueryParam("category", category)
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	g.apiTest.SendRequest()
	return nil
}

func (g *getPermissionsTest) iShouldGetAllPermissionsInThatCategory() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var permissionsList []permissions.Permission
	err := g.apiTest.UnmarshalResponseBodyPath("data", &permissionsList)
	if err != nil {
		return err
	}

	query := "select v0,v1,v2 from casbin_rule where p_type = 'p'"
	if g.category != "" {
		query = fmt.Sprintf("%s and v2 = '%s'", query, g.category)
	}
	var dbPermissions []permissions.Permission
	rows, err := g.Conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	for rows.Next() {
		var i permissions.Permission
		if err := rows.Scan(&i.ID, &i.Name, &i.Category); err != nil {
			return err
		}
		dbPermissions = append(dbPermissions, i)
	}

	for _, v := range permissionsList {
		if !collection.ContainsWithMatcher(v.ID, dbPermissions, func(value string, perm permissions.Permission) bool {
			return value == perm.ID
		}) {
			return fmt.Errorf("expected to get: %v", v)
		}
	}
	if g.category != "" {
		for _, v := range permissionsList {
			if err := g.apiTest.AssertEqual(v.Category, g.category); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *getPermissionsTest) myRequestShouldFailWithMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if err := g.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (g *getPermissionsTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.Admin.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, g.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I request to get all permissions with category "([^"]*)"$`, g.iRequestToGetAllPermissionsWithCategory)
	ctx.Step(`^I should get all permissions in that category$`, g.iShouldGetAllPermissionsInThatCategory)
	ctx.Step(`^my request should fail with message "([^"]*)"$`, g.myRequestShouldFailWithMessage)
}
