package create_resource_server

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
)

type createResourceServerTest struct {
	test.TestInstance
	apiTest        src.ApiTest
	resourceServer dto.ResourceServer
	existingServer db.ResourceServer
	Admin          db.User
}

func TestCreateResourceServer(t *testing.T) {
	r := &createResourceServerTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest = src.ApiTest{
		Server: r.Server,
	}
	r.apiTest.URL = "/v1/resourceServer"
	r.apiTest.Method = http.MethodPost
	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.InitializeTest(t, "create resource server test", "features/create_resource_server.feature", r.InitializeScenario)
}

// background functions
func (r *createResourceServerTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error
	r.Admin, err = r.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	return r.GrantRoleForUser(r.Admin.ID.String(), adminCredentials)
}

// Given functions
func (r *createResourceServerTest) iHaveFilledResourceServerNameAndTheFollowingScopes(serverName string, scopesTable *godog.Table) error {
	scopes, err := r.apiTest.ReadRows(scopesTable, nil, false)
	if err != nil {
		return err
	}
	r.apiTest.SetBodyMap(map[string]interface{}{
		"name":   serverName,
		"scopes": scopes,
	})

	r.resourceServer = dto.ResourceServer{
		Name: serverName,
	}
	if err := r.apiTest.UnmarshalJSON([]byte(scopes), &r.resourceServer.Scopes); err != nil {
		return err
	}

	return nil
}
func (r *createResourceServerTest) theResourceServerIsRegistered(serverName string) error {
	resourceServer, err := r.DB.CreateResourceServer(context.Background(), serverName)
	if err != nil {
		return err
	}
	r.existingServer = resourceServer

	return nil
}

// When functions
func (r *createResourceServerTest) iSubmitToCreateAResourceServer() error {
	r.apiTest.SetHeader("Authorization", "Bearer "+r.AccessToken)
	r.apiTest.SendRequest()

	return nil
}

// Then functions
func (r *createResourceServerTest) theResourceServerShouldBeCreated() error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var resourceServer dto.ResourceServer
	if err := r.apiTest.UnmarshalResponseBodyPath("data", &resourceServer); err != nil {
		return err
	}
	if err := r.apiTest.AssertEqual(resourceServer.Name, r.resourceServer.Name); err != nil {
		return err
	}
	for _, v := range resourceServer.Scopes {
		found := false
		for k, v2 := range r.resourceServer.Scopes {
			if v.Name == fmt.Sprintf(r.resourceServer.Name, ".", v2.Name) {
				r.resourceServer.Scopes[k].Name = v.Name // to hold the real scope name for deleting it after the test
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected scope: %v", v)
		}
	}

	return nil
}
func (r *createResourceServerTest) theRequestShouldFailWithAnd(message, fieldError string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if message != "" {
		if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
			return err
		}
	}
	if fieldError != "" {
		if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", fieldError); err != nil {
			return err
		}
	}

	return nil
}

func (r *createResourceServerTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(ctx, r.Admin.ID)
		_, _ = r.DB.DeleteResourceServer(ctx, r.resourceServer.ID)
		_, _ = r.DB.DeleteResourceServer(ctx, r.existingServer.ID)
		for _, v := range r.resourceServer.Scopes {
			_, _ = r.DB.DeleteScope(ctx, v.Name)
		}
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, r.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I have filled resource server name "([^"]*)" and the following scopes$`, r.iHaveFilledResourceServerNameAndTheFollowingScopes)
	ctx.Step(`^I submit to create a resource server$`, r.iSubmitToCreateAResourceServer)
	ctx.Step(`^the request should fail with "([^"]*)" and "([^"]*)"$`, r.theRequestShouldFailWithAnd)
	ctx.Step(`^the resource server "([^"]*)" is registered$`, r.theResourceServerIsRegistered)
	ctx.Step(`^the resource server should be created$`, r.theResourceServerShouldBeCreated)
}
