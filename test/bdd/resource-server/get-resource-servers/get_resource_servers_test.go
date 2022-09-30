package get_resource_servers

import (
	"context"
	"database/sql"
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

type getResourceServersTest struct {
	test.TestInstance
	apiTest         src.ApiTest
	resourceServers []db.ResourceServer
	scopes          []db.Scope
	Admin           db.User
	Preferences     preferenceData
}

type preferenceData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func TestGetClients(t *testing.T) {
	c := getResourceServersTest{}
	c.apiTest.URL = "/v1/resourceServers"
	c.apiTest.Method = http.MethodGet
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "get resource servers test", "features/get_resource_servers.feature", c.InitializeScenario)
}

func (c *getResourceServersTest) theFollowingResourceServersAreRegisteredOnTheSystem(resourceServers *godog.Table) error {
	resourceServersJSON, err := c.apiTest.ReadRows(resourceServers, nil, false)
	if err != nil {
		return err
	}
	var resourceServersData []db.ResourceServer
	err = c.apiTest.UnmarshalJSON([]byte(resourceServersJSON), &resourceServersData)
	if err != nil {
		return err
	}
	for _, v := range resourceServersData {
		resourceServer, err := c.DB.CreateResourceServer(context.Background(), v.Name)
		if err != nil {
			return err
		}
		c.resourceServers = append(c.resourceServers, resourceServer)
	}

	return nil
}
func (c *getResourceServersTest) theResourceServersHaveTheFollowingScopes(scopesTable *godog.Table) error {
	scopesJSON, err := c.apiTest.ReadRows(scopesTable, nil, false)
	if err != nil {
		return err
	}
	var scopesData []dto.Scope
	err = c.apiTest.UnmarshalJSON([]byte(scopesJSON), &scopesData)
	if err != nil {
		return err
	}
	for _, v := range scopesData {
		scope, err := c.DB.CreateScope(context.Background(), db.CreateScopeParams{
			Name:        v.Name,
			Description: v.Description,
			ResourceServerName: sql.NullString{
				String: v.ResourceServerName,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
		c.scopes = append(c.scopes, scope)
	}
	return nil
}
func (c *getResourceServersTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	c.Admin, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	return c.GrantRoleForUser(c.Admin.ID.String(), adminCredentials)
}

func (c *getResourceServersTest) iRequestToGetAllTheResourceServersWithTheFollowingPreferences(preferences *godog.Table) error {
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

func (c *getResourceServersTest) iShouldGetTheListOfResourceServersThatPassMyPreferences() error {
	var responseResourceServers []dto.ResourceServer
	var metaData model.MetaData

	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	err := c.apiTest.UnmarshalResponseBodyPath("meta_data", &metaData)
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalResponseBodyPath("data", &responseResourceServers)
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
	if err := c.apiTest.AssertEqual(len(responseResourceServers), total); err != nil {
		return err
	}
	for _, v := range responseResourceServers {
		found := false
		for _, v2 := range c.resourceServers {
			if v.ID.String() == v2.ID.String() {
				found = true
				continue
			}
		}
		if !found {
			return fmt.Errorf("expected resource server: %v", v)
		}
	}
	return nil
}

func (c *getResourceServersTest) iShouldGetErrorMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (c *getResourceServersTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for _, v := range c.resourceServers {
			_, _ = c.DB.DeleteResourceServer(ctx, v.ID)
		}
		for _, v := range c.scopes {
			_, _ = c.DB.DeleteScope(ctx, v.Name)
		}
		_, _ = c.DB.DeleteUser(ctx, c.Admin.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in as admin user$`, c.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get all the resource servers with the following preferences$`, c.iRequestToGetAllTheResourceServersWithTheFollowingPreferences)
	ctx.Step(`^I should get the list of resource servers that pass my preferences$`, c.iShouldGetTheListOfResourceServersThatPassMyPreferences)
	ctx.Step(`^I should get error message "([^"]*)"$`, c.iShouldGetErrorMessage)
	ctx.Step(`^The following resource servers are registered on the system$`, c.theFollowingResourceServersAreRegisteredOnTheSystem)
	ctx.Step(`^the resource servers have the following scopes$`, c.theResourceServersHaveTheFollowingScopes)

}
