package update_status

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

type updateClientStatusTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Admin   db.User
	client  db.Client
}

func TestUpdateClientStatus(t *testing.T) {
	u := updateClientStatusTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.InitializeTest(t, "update client status", "features/update_status.feature", u.InitializeScenario)
}

func (u *updateClientStatusTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
	body, err := u.apiTest.ReadRow(adminCredential, nil, false)
	if err != nil {
		return err
	}

	adminValue := dto.User{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &adminValue)
	if err != nil {
		return err
	}

	u.Admin, err = u.AuthenticateWithParam(adminValue)
	if err != nil {
		return err
	}
	u.apiTest.SetHeader("Authorization", "Bearer "+u.AccessToken)
	return u.GrantRoleForUser(u.Admin.ID.String(), adminCredential)
}

func (u *updateClientStatusTest) thereIsClientWithTheFollowingDetails(client *godog.Table) error {

	clientJSON, err := u.apiTest.ReadRow(client, nil, false)
	if err != nil {
		return err
	}

	var clientData db.CreateClientParams
	err = u.apiTest.UnmarshalJSON([]byte(clientJSON), &clientData)
	if err != nil {
		return err
	}

	u.client, err = u.DB.CreateClient(context.Background(), clientData)
	if err != nil {
		return err
	}

	u.apiTest.URL += u.client.ID.String() + "/status"

	return nil
}

func (u *updateClientStatusTest) iUpdateTheClientsStatusTo(updatedStatus string) error {
	u.apiTest.SetBodyMap(map[string]interface{}{
		"status": updatedStatus,
	})
	u.apiTest.SendRequest()
	return nil
}

func (u *updateClientStatusTest) theClientStatusShouldUpdateTo(updatedStatus string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	updatedClientData, err := u.DB.GetClientByID(context.Background(), u.client.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(updatedClientData.Status, updatedStatus); err != nil {
		return err
	}
	return nil
}

func (u *updateClientStatusTest) thenIShouldGetClientNotFoundErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (u *updateClientStatusTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (u *updateClientStatusTest) thereIsClientWithId(clientID string) error {
	u.apiTest.URL += clientID + "/status"
	return nil
}

func (u *updateClientStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {
	u.apiTest.URL = "/v1/clients/"
	u.apiTest.Method = http.MethodPatch
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteClient(ctx, u.client.ID)
		_, _ = u.DB.DeleteUser(ctx, u.Admin.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, u.iAmLoggedInAsAdminUser)
	ctx.Step(`^I update the client\'s status to "([^"]*)"$`, u.iUpdateTheClientsStatusTo)
	ctx.Step(`^the client status should update to "([^"]*)"$`, u.theClientStatusShouldUpdateTo)
	ctx.Step(`^Then I should get client not found error with message "([^"]*)"$`, u.thenIShouldGetClientNotFoundErrorWithMessage)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, u.thenIShouldGetErrorWithMessage)
	ctx.Step(`^there is client with id "([^"]*)"$`, u.thereIsClientWithId)
	ctx.Step(`^there is client with the following details:$`, u.thereIsClientWithTheFollowingDetails)
}
