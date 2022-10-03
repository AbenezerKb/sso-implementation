package get_client

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

type getClientTest struct {
	test.TestInstance
	apiTest src.ApiTest
	client  db.Client
	Admin   db.User
}

func TestGetClient(t *testing.T) {
	g := getClientTest{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeTest(t, "get client test", "features/get_client.feature", g.InitializeScenario)
}

func (g *getClientTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	g.Admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	return g.GrantRoleForUser(g.Admin.ID.String(), adminCredentials)
}

func (g *getClientTest) thereIsClientWithTheFollowingDetails(clientDetails *godog.Table) error {
	body, err := g.apiTest.ReadRow(clientDetails, nil, false)
	if err != nil {
		return err
	}

	clientValues := db.CreateClientParams{}
	err = g.apiTest.UnmarshalJSON([]byte(body), &clientValues)
	if err != nil {
		return err
	}

	g.client, err = g.DB.CreateClient(context.Background(), clientValues)
	if err != nil {
		return err
	}

	return nil
}
func (g *getClientTest) iGetTheClient() error {
	g.apiTest.SendRequest()

	return nil
}

func (g *getClientTest) iHaveClientId() error {
	g.apiTest.URL += g.client.ID.String()

	return nil
}

func (g *getClientTest) iHaveClientWithId(id string) error {
	g.apiTest.URL += id

	return nil
}

func (g *getClientTest) iShouldSuccessfullyGetTheClient() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	respondedClient := dto.Client{}
	err := g.apiTest.UnmarshalResponseBodyPath("data", &respondedClient)
	if err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(respondedClient.ID, g.client.ID); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(respondedClient.Name, g.client.Name); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(respondedClient.ClientType, g.client.ClientType); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(respondedClient.LogoURL, g.client.LogoUrl); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(respondedClient.Status, g.client.Status); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(respondedClient.Scopes, g.client.Scopes); err != nil {
		return err
	}

	return nil
}

func (g *getClientTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return g.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (g *getClientTest) InitializeScenario(ctx *godog.ScenarioContext) {
	g.apiTest.URL = "/v1/clients/"
	g.apiTest.Method = http.MethodGet
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.apiTest.InitializeServer(g.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.Admin.ID)
		_, _ = g.DB.DeleteClient(ctx, g.client.ID)

		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I Get the client$`, g.iGetTheClient)
	ctx.Step(`^I have client id$`, g.iHaveClientId)
	ctx.Step(`^I have client with id "([^"]*)"$`, g.iHaveClientWithId)
	ctx.Step(`^I should successfully get the client$`, g.iShouldSuccessfullyGetTheClient)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, g.thenIShouldGetErrorWithMessage)
	ctx.Step(`^there is client with the following details:$`, g.thereIsClientWithTheFollowingDetails)
}
