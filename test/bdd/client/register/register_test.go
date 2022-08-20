package register

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

type clientRegistrationTest struct {
	test.TestInstance
	apiTest src.ApiTest
	client  *dto.Client
	Admin   db.User
}

func TestClientRegistration(t *testing.T) {
	c := &clientRegistrationTest{}
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeTest(t, "Client registration test", "features/client_registration.feature", c.InitializeScenario)
}
func (c *clientRegistrationTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	c.Admin, err = c.Authenicate(adminCredentials)
	if err != nil {
		return err
	}
	return c.GrantRoleForUser(c.Admin.ID.String(), adminCredentials)
}

func (c *clientRegistrationTest) iFillTheFollowingClientForm(clientForm *godog.Table) error {
	body, err := c.apiTest.ReadRow(clientForm, []src.Type{
		{
			Column: "redirect_uris",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	c.apiTest.Body = body
	return nil
}

func (c *clientRegistrationTest) iSubmitTheForm() error {
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	c.apiTest.SendRequest()
	return nil
}

func (c *clientRegistrationTest) theRegistrationShouldBeSuccessful() error {
	if err := c.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	return c.apiTest.UnmarshalResponseBodyPath("data", &c.client)
}

func (c *clientRegistrationTest) theRegistrationShouldFailWith(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (c *clientRegistrationTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/clients"
		c.apiTest.Method = http.MethodPost
		c.apiTest.SetHeader("Content-Type", "application/json")
		c.apiTest.InitializeServer(c.Server)
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {

		// delete the registered client
		//_, err = c.DB.DeleteClient(ctx, c.user.Data.ID)
		//if err != nil {
		//	return ctx, err
		//}

		// delete the admin
		_, err = c.DB.DeleteUser(ctx, c.Admin.ID)
		return ctx, err
	})
	ctx.Step(`^I am logged in as admin user$`, c.iAmLoggedInAsAdminUser)
	ctx.Step(`^I fill the following client form$`, c.iFillTheFollowingClientForm)
	ctx.Step(`^I submit the form$`, c.iSubmitTheForm)
	ctx.Step(`^The registration should be successful$`, c.theRegistrationShouldBeSuccessful)
	ctx.Step(`^The registration should fail with "([^"]*)"$`, c.theRegistrationShouldFailWith)
}
