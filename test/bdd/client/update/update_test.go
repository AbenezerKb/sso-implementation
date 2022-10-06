package update

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

type updateClientTest struct {
	test.TestInstance
	apiTest      src.ApiTest
	Admin        db.User
	client       db.Client
	updateClient dto.Client
}

func TestUpdateUserStatus(t *testing.T) {
	u := updateClientTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.InitializeTest(t, "update client", "features/update.feature", u.InitializeScenario)
}

func (u *updateClientTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
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

func (u *updateClientTest) thereIsClientWithTheFollowingDetails(clientDetails *godog.Table) error {
	body, err := u.apiTest.ReadRow(clientDetails, nil, false)

	if err != nil {
		return err
	}

	clientValues := db.CreateClientParams{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &clientValues)
	if err != nil {
		return err
	}
	u.client, err = u.DB.CreateClient(context.Background(), clientValues)
	u.apiTest.URL += u.client.ID.String()

	return err
}

func (u *updateClientTest) iFillTheFormWithTheFollowingDetails(clientForm *godog.Table) error {
	body, err := u.apiTest.ReadRow(clientForm, []src.Type{
		{
			Column: "redirect_uris",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}

	clientValues := dto.Client{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &clientValues)
	if err != nil {
		return err
	}

	u.updateClient = clientValues

	// jsonClient, err := json.Marshal(clientValues)
	// if err != nil {
	// 	return err
	// }

	// u.apiTest.Body = string(jsonClient)
	u.apiTest.Body = body
	return nil
}

func (u *updateClientTest) iUpdateTheClient() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *updateClientTest) theClientShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	client, err := u.DB.GetClientByID(context.Background(), u.client.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(u.client.Name, client.Name); err != nil {
		return err
	}
	if err := u.apiTest.AssertEqual(u.client.ClientType, client.ClientType); err != nil {
		return err
	}
	if err := u.apiTest.AssertEqual(u.client.LogoUrl, client.LogoUrl); err != nil {
		return err
	}
	if err := u.apiTest.AssertEqual(u.client.Scopes, client.Scopes); err != nil {
		return err
	}
	if err := u.apiTest.AssertEqual(u.client.RedirectUris, client.RedirectUris); err != nil {
		return err
	}

	return nil
}

func (u *updateClientTest) theClientUpdatedShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (u *updateClientTest) InitializeScenario(ctx *godog.ScenarioContext) {

	u.apiTest.URL = "/v1/clients/"
	u.apiTest.Method = http.MethodPut
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteClient(ctx, u.client.ID)
		_, _ = u.DB.DeleteUser(ctx, u.Admin.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, u.iAmLoggedInAsAdminUser)
	ctx.Step(`^I fill the form with the following details$`, u.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I update the client$`, u.iUpdateTheClient)
	ctx.Step(`^The client should be updated$`, u.theClientShouldBeUpdated)
	ctx.Step(`^The client updated should fail with message "([^"]*)"$`, u.theClientUpdatedShouldFailWithMessage)
	ctx.Step(`^There is client with the following details$`, u.thereIsClientWithTheFollowingDetails)
}
