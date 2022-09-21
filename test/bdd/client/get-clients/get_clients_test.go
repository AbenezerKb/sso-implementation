package get_clients

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
)

type getClientsTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	clients     []db.Client
	Admin       db.User
	Preferences preferenceData
}

type preferenceData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func TestGetClients(t *testing.T) {
	c := getClientsTest{}
	c.apiTest.URL = "/v1/clients"
	c.apiTest.Method = http.MethodGet
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "get clients test", "features/get_clients.feature", c.InitializeScenario)
}

func (c *getClientsTest) theFollowingClientsAreRegisteredOnTheSystem(clients *godog.Table) error {
	clientsJSON, err := c.apiTest.ReadRows(clients, nil, false)
	if err != nil {
		return err
	}
	var clientsData []db.CreateClientParams
	err = c.apiTest.UnmarshalJSON([]byte(clientsJSON), &clientsData)
	if err != nil {
		return err
	}
	for _, v := range clientsData {
		client, err := c.DB.CreateClient(context.Background(), v)
		if err != nil {
			return err
		}
		c.clients = append(c.clients, client)
	}

	return nil
}

func (c *getClientsTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	c.Admin, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	return c.GrantRoleForUser(c.Admin.ID.String(), adminCredentials)
}

func (c *getClientsTest) iRequestToGetAllTheClientsWithTheFollowingPreferences(preferences *godog.Table) error {
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

func (c *getClientsTest) iShouldGetTheListOfClientsThatPassMyPreferences() error {
	var responseClients []dto.Client
	var metaData model.MetaData

	err := c.apiTest.UnmarshalResponseBodyPath("meta_data", &metaData)
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalResponseBodyPath("data", &responseClients)
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
	if err := c.apiTest.AssertEqual(len(responseClients), total); err != nil {
		return err
	}
	for _, v := range responseClients {
		found := false
		for _, v2 := range c.clients {
			if v.ID.String() == v2.ID.String() {
				found = true
				continue
			}
		}
		if !found {
			return fmt.Errorf("expected client: %v", v)
		}
	}
	return nil
}

func (c *getClientsTest) iShouldGetErrorMessage(message string) error {
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (c *getClientsTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for _, v := range c.clients {
			_, _ = c.DB.DeleteClient(ctx, v.ID)
		}
		_, _ = c.DB.DeleteUser(ctx, c.Admin.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in as admin user$`, c.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get all the clients with the following preferences$`, c.iRequestToGetAllTheClientsWithTheFollowingPreferences)
	ctx.Step(`^I should get the list of clients that pass my preferences$`, c.iShouldGetTheListOfClientsThatPassMyPreferences)
	ctx.Step(`^I should get error message "([^"]*)"$`, c.iShouldGetErrorMessage)
	ctx.Step(`^The following clients are registered on the system$`, c.theFollowingClientsAreRegisteredOnTheSystem)
}
