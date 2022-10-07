package get_roles

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils/collection"
	"sso/test"
	"testing"
)

type getRolesTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	roles       []dto.Role
	Admin       db.User
	adminRole   *dto.Role
	Preferences preferenceData
}

type preferenceData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func TestGetRoles(t *testing.T) {
	r := getRolesTest{}
	r.apiTest.URL = "/v1/roles"
	r.apiTest.Method = http.MethodGet
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "get roles test", "features/get_roles.feature", r.InitializeScenario)
}

func (c *getRolesTest) theFollowingRolesAreRegisteredOnTheSystem(roles *godog.Table) error {
	rolesJSON, err := c.apiTest.ReadRows(roles, []src.Type{
		{
			Column: "permissions",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	var rolesData []dto.Role
	err = c.apiTest.UnmarshalJSON([]byte(rolesJSON), &rolesData)
	if err != nil {
		return err
	}
	for _, v := range rolesData {
		role, err := c.PersistDB.CreateRoleTX(context.Background(), v.Name, v.Permissions)
		if err != nil {
			return err
		}
		c.roles = append(c.roles, role)
	}

	return nil
}

func (c *getRolesTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	c.Admin, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	c.adminRole, c.GrantRoleAfterFunc, err = c.GrantRoleForUserWithAfter(c.Admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}
	return nil
}

func (c *getRolesTest) iRequestToGetAllTheRolesWithTheFollowingPreferences(preferences *godog.Table) error {
	preferencesJSON, err := c.apiTest.ReadRow(preferences, []src.Type{
		{
			Column: "page",
			Kind:   src.Any,
		},
		{
			Column: "per_page",
			Kind:   src.Any,
		},
	}, false)
	if err != nil {
		return err
	}
	err = c.apiTest.UnmarshalJSON([]byte(preferencesJSON), &c.Preferences)
	if err != nil {
		return err
	}

	c.apiTest.SetQueryParam("page", fmt.Sprintf("%d", c.Preferences.Page))
	c.apiTest.SetQueryParam("per_page", fmt.Sprintf("%d", c.Preferences.PerPage))
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	c.apiTest.SendRequest()
	return nil
}

func (c *getRolesTest) iShouldGetTheListOfRolesThatPassMyPreferences() error {
	var responseRoles []dto.Role
	var metaData model.MetaData

	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	err := c.apiTest.UnmarshalResponseBodyPath("meta_data", &metaData)
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalResponseBodyPath("data", &responseRoles)
	if err != nil {
		return err
	}
	var total int
	if c.Preferences.Page < metaData.Total/c.Preferences.PerPage {
		total = c.Preferences.PerPage
	} else if c.Preferences.Page == metaData.Total/c.Preferences.PerPage {
		total = metaData.Total % c.Preferences.PerPage
	} else {
		total = 0
	}
	if err := c.apiTest.AssertEqual(len(responseRoles), total); err != nil {
		return err
	}
	for _, v := range responseRoles {
		found := false
		for _, v2 := range append(c.roles, *c.adminRole) {
			if v.Name == v2.Name {
				found = true
				for _, permission := range v.Permissions {
					if !collection.Contains(permission, v2.Permissions) {
						return fmt.Errorf("expected permission `%s` in role %s", permission, v2.Name)
					}
				}
				continue
			}
		}
		if !found {
			return fmt.Errorf("expected role: %v", v)
		}
	}
	return nil
}

func (c *getRolesTest) iShouldGetErrorMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (c *getRolesTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// FIXME: this is not correct but is used to insure compatibility with the legacy function `GrantRoleForUser`
		_, _ = c.Conn.Exec(ctx, "DELETE FROM roles WHERE true")
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for _, v := range c.roles {
			_, _ = c.DB.DeleteRole(ctx, v.Name)
			_, _ = c.Conn.Exec(ctx, "DELETE FROM casbin_rule WHERE v0 = $1", v.Name)
		}
		_ = c.GrantRoleAfterFunc()
		_, _ = c.DB.DeleteUser(ctx, c.Admin.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in as admin user$`, c.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get all the roles with the following preferences$`, c.iRequestToGetAllTheRolesWithTheFollowingPreferences)
	ctx.Step(`^I should get the list of roles that pass my preferences$`, c.iShouldGetTheListOfRolesThatPassMyPreferences)
	ctx.Step(`^I should get error message "([^"]*)"$`, c.iShouldGetErrorMessage)
	ctx.Step(`^The following roles are registered on the system$`, c.theFollowingRolesAreRegisteredOnTheSystem)
}
