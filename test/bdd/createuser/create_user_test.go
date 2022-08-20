package createuser

import (
	"context"
	"encoding/json"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type createuserTest struct {
	test.TestInstance
	apiTest src.ApiTest
	user    struct {
		OK   bool     `json:"ok"`
		Data dto.User `json:"data"`
	}
	Admin db.User
}

func TestCreateuser(t *testing.T) {

	c := &createuserTest{}
	c.TestInstance = test.Initiate("../../../")

	c.apiTest.InitializeTest(t, "Create user test", "features/create_user.feature", c.InitializeScenario)
}
func (c *createuserTest) iAmLoggedInWithTheFollowingCreadentials(creadentials *godog.Table) error {
	var err error
	c.Admin, err = c.Authenicate(creadentials)
	if err != nil {
		return err
	}
	return c.GrantRoleForUser(c.Admin.ID.String(), creadentials)
}

func (c *createuserTest) iFillTheFormWithTheFollowingDetails(userForm *godog.Table) error {
	body, err := c.apiTest.ReadRow(userForm, nil, false)
	if err != nil {
		return err
	}
	c.apiTest.Body = body
	return nil
}

func (c *createuserTest) iSubmitTheCreateUserForm() error {
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	c.apiTest.SendRequest()
	return nil
}

func (c *createuserTest) theCreatingProcessShouldFailWith(msg string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", msg); err != nil {
		return err
	}
	return nil

}

func (c *createuserTest) theUserIsCreated() error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	err := json.Unmarshal(c.apiTest.ResponseBody, &c.user)

	return err
}

func (c *createuserTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/users"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		c.apiTest.InitializeServer(c.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// delete the registered user
		_, _ = c.DB.DeleteUser(ctx, c.user.Data.ID)

		// delete the admin
		_, _ = c.DB.DeleteUser(ctx, c.Admin.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in with the following creadentials$`, c.iAmLoggedInWithTheFollowingCreadentials)
	ctx.Step(`^I fill the form with the following details$`, c.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I submit the create user form$`, c.iSubmitTheCreateUserForm)
	ctx.Step(`^the creating process should fail with "([^"]*)"$`, c.theCreatingProcessShouldFailWith)
	ctx.Step(`^The user is created$`, c.theUserIsCreated)
}
