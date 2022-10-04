package create_role

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

type createRoleTest struct {
	test.TestInstance
	apiTest src.ApiTest
	role    dto.Role
	admin   db.User
}

func TestCreateRole(t *testing.T) {
	c := &createRoleTest{}
	c.apiTest.URL = "/v1/roles"
	c.apiTest.Method = http.MethodPost
	c.apiTest.SetHeader("Content-Type", "application/json")
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "create role test", "features/create_role.feature", c.InitializeScenario)
}

func (c *createRoleTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error
	c.admin, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	return c.GrantRoleForUser(c.admin.ID.String(), adminCredentials)
}

func (c *createRoleTest) iRequestToCreateARoleWithTheFollowingPermissions(roleTable *godog.Table) error {
	body, err := c.apiTest.ReadRow(roleTable, []src.Type{
		{
			Column:   "role_name",
			WithName: "name",
		},
		{
			Column: "permissions",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	err = c.apiTest.UnmarshalJSON([]byte(body), &c.role)
	if err != nil {
		return err
	}
	c.apiTest.Body = body
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	c.apiTest.SendRequest()
	return nil
}

func (c *createRoleTest) myRequestShouldFailWithAnd(message, fieldError string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if message != "" {
		if err := c.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
			return err
		}
	}
	if fieldError != "" {
		if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_errors.0.description", fieldError); err != nil {
			return err
		}
	}
	return nil
}

func (c *createRoleTest) theRoleShouldSuccessfullyBeCreated() error {
	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	var role dto.Role
	err := c.apiTest.UnmarshalResponseBodyPath("data", &role)
	if err != nil {
		return err
	}

	if err := c.apiTest.AssertEqual(role.Name, c.role.Name); err != nil {
		return err
	}
	for i := 0; i < len(role.Permissions); i++ {
		if !collection.Contains(role.Permissions[i], c.role.Permissions) {
			return fmt.Errorf("expected role: %v", role)
		}
	}

	return nil
}

func (c *createRoleTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		_, _ = c.DB.DeleteUser(ctx, c.admin.ID)
		// TODO: delete role here
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, c.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I request to create a role with the following permissions$`, c.iRequestToCreateARoleWithTheFollowingPermissions)
	ctx.Step(`^my request should fail with "([^"]*)" and "([^"]*)"$`, c.myRequestShouldFailWithAnd)
	ctx.Step(`^the role should successfully be created$`, c.theRoleShouldSuccessfullyBeCreated)
}
