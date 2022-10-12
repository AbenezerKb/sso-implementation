package update_identity_provider

import (
	"context"
	"database/sql"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateIdentityProviderTest struct {
	test.TestInstance
	apiTest                 src.ApiTest
	identityProvider        db.IdentityProvider
	updatedIdentityProvider dto.IdentityProvider
	Admin                   db.User
}

func TestUpdateIdentityProvider(t *testing.T) {
	u := &updateIdentityProviderTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest = src.ApiTest{
		Server: u.Server,
	}

	u.apiTest.InitializeTest(t, "update identity provider test", "features/update_identity_provider.feature", u.InitializeScenario)
}

func (u *updateIdentityProviderTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	u.Admin, err = u.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, u.GrantRoleAfterFunc, err = u.GrantRoleForUserWithAfter(u.Admin.ID.String(), adminCredentials)
	u.apiTest.SetHeader("Authorization", "Bearer "+u.AccessToken)

	return err
}

func (u *updateIdentityProviderTest) thereIsIdentityProviderWithTheFollowingDetails(idpTable *godog.Table) error {
	idpJSON, err := u.apiTest.ReadRow(idpTable, nil, false)
	if err != nil {
		return err
	}

	idpData := dto.IdentityProvider{}
	if err := u.apiTest.UnmarshalJSON([]byte(idpJSON), &idpData); err != nil {
		return err
	}

	u.identityProvider, err = u.DB.CreateIdentityProvider(context.Background(), db.CreateIdentityProviderParams{
		Name:                idpData.Name,
		LogoUrl:             sql.NullString{String: idpData.LogoURI, Valid: true},
		ClientSecret:        idpData.ClientSecret,
		ClientID:            idpData.ClientID,
		RedirectUri:         idpData.RedirectURI,
		AuthorizationUri:    idpData.AuthorizationURI,
		TokenEndpointUrl:    idpData.TokenEndpointURI,
		UserInfoEndpointUrl: sql.NullString{String: idpData.UserInfoEndpointURI, Valid: true},
	})
	if err != nil {
		return err
	}

	u.apiTest.URL += u.identityProvider.ID.String()

	return nil
}

func (u *updateIdentityProviderTest) iFillTheFormWithTheFollowingDetails(idpTable *godog.Table) error {
	ipJSON, err := u.apiTest.ReadRow(idpTable, nil, false)
	if err != nil {
		return err
	}
	u.apiTest.Body = ipJSON
	return u.apiTest.UnmarshalJSON([]byte(ipJSON), &u.updatedIdentityProvider)
}

func (u *updateIdentityProviderTest) iUpdateTheIdentityProvider() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *updateIdentityProviderTest) theIdentityProviderShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	fetchedIdP, err := u.DB.GetIdentityProvider(context.Background(), u.identityProvider.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.ClientID, u.updatedIdentityProvider.ClientID); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.ClientSecret, u.updatedIdentityProvider.ClientSecret); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.LogoUrl.String, u.updatedIdentityProvider.LogoURI); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.TokenEndpointUrl, u.updatedIdentityProvider.TokenEndpointURI); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.UserInfoEndpointUrl.String, u.updatedIdentityProvider.UserInfoEndpointURI); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.Name, u.updatedIdentityProvider.Name); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.AuthorizationUri, u.updatedIdentityProvider.AuthorizationURI); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedIdP.RedirectUri, u.updatedIdentityProvider.RedirectURI); err != nil {
		return err
	}

	return nil
}

func (u *updateIdentityProviderTest) theIdentityProviderUpdatedShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (u *updateIdentityProviderTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		u.apiTest.URL = "/v1/identityProviders/"
		u.apiTest.Method = http.MethodPut
		u.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteUser(ctx, u.Admin.ID)
		_, _ = u.DB.DeleteIdentityProvider(ctx, u.identityProvider.ID)
		_ = u.GrantRoleAfterFunc()
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, u.iAmLoggedInAsAdminUser)
	ctx.Step(`^I fill the form with the following details$`, u.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I update the identity provider$`, u.iUpdateTheIdentityProvider)
	ctx.Step(`^The identity provider should be updated$`, u.theIdentityProviderShouldBeUpdated)
	ctx.Step(`^The identity provider updated should fail with message "([^"]*)"$`, u.theIdentityProviderUpdatedShouldFailWithMessage)
	ctx.Step(`^There is identity provider with the following details$`, u.thereIsIdentityProviderWithTheFollowingDetails)
}
