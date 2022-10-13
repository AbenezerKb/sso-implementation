package delete

import (
	"context"
	"fmt"
	"net/http"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type deleteClientTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Admin   db.User
	client  db.Client
}

func TestDeleteClient(t *testing.T) {
	d := deleteClientTest{}
	d.TestInstance = test.Initiate("../../../../")
	d.apiTest.InitializeTest(t, "delete client", "features/delete_client.feature", d.InitializeScenario)
}
func (d *deleteClientTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
	body, err := d.apiTest.ReadRow(adminCredential, nil, false)
	if err != nil {
		return err
	}

	adminValue := dto.User{}
	err = d.apiTest.UnmarshalJSON([]byte(body), &adminValue)
	if err != nil {
		return err
	}

	d.Admin, err = d.AuthenticateWithParam(adminValue)
	if err != nil {
		return err
	}
	d.apiTest.SetHeader("Authorization", "Bearer "+d.AccessToken)
	return d.GrantRoleForUser(d.Admin.ID.String(), adminCredential)
}
func (d *deleteClientTest) thereIsAClientWithTheFollowingDetails(clientDetails *godog.Table) error {
	body, err := d.apiTest.ReadRow(clientDetails, nil, false)

	if err != nil {
		return err
	}

	clientValues := db.CreateClientParams{}
	err = d.apiTest.UnmarshalJSON([]byte(body), &clientValues)
	if err != nil {
		return err
	}
	d.client, err = d.DB.CreateClient(context.Background(), clientValues)

	return err
}

func (d *deleteClientTest) iDeleteTheClient() error {
	d.apiTest.URL += d.client.ID.String()
	d.apiTest.SendRequest()
	return nil
}

func (d *deleteClientTest) iDeleteTheClientWithId(clientId string) error {
	d.apiTest.URL += clientId
	d.apiTest.SendRequest()
	return nil
}

func (d *deleteClientTest) theClientShouldBeDeleted() error {
	if err := d.apiTest.AssertStatusCode(http.StatusNoContent); err != nil {
		return err
	}

	_, err := d.DB.GetClientByID(context.Background(), d.client.ID)
	if err == nil {
		return fmt.Errorf("client is not deleted")
	}
	if !sqlcerr.Is(err, sqlcerr.ErrNoRows) {
		return err
	}

	return nil
}

func (d *deleteClientTest) theDeleteShouldFailWithErrorMessage(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return d.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}
func (d *deleteClientTest) InitializeScenario(ctx *godog.ScenarioContext) {
	d.apiTest.URL = "/v1/clients/"
	d.apiTest.Method = http.MethodDelete
	d.apiTest.SetHeader("Content-Type", "application/json")
	d.apiTest.InitializeServer(d.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.DeleteUser(ctx, d.Admin.ID)
		_, _ = d.DB.DeleteClient(ctx, d.client.ID)

		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, d.iAmLoggedInAsAdminUser)
	ctx.Step(`^I delete the client$`, d.iDeleteTheClient)
	ctx.Step(`^The client should be deleted$`, d.theClientShouldBeDeleted)
	ctx.Step(`^There is a client with the following details$`, d.thereIsAClientWithTheFollowingDetails)
	ctx.Step(`^I delete the client with id "([^"]*)"$`, d.iDeleteTheClientWithId)
	ctx.Step(`^The delete should fail with error message "([^"]*)"$`, d.theDeleteShouldFailWithErrorMessage)

}
