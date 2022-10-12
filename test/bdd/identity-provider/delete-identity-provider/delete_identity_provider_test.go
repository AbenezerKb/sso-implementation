package main

import (
	"context"
	"database/sql"
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

type deleteIdentityProviderTest struct {
	test.TestInstance
	apiTest          src.ApiTest
	identityProvider db.IdentityProvider
	Admin            db.User
}

func TestUpdateIdentityProvider(t *testing.T) {
	d := &deleteIdentityProviderTest{}
	d.TestInstance = test.Initiate("../../../../")
	d.apiTest = src.ApiTest{
		Server: d.Server,
	}

	d.apiTest.InitializeTest(t, "delete identity provider test", "features/delete_identity_provider.feature", d.InitializeScenario)
}

func (d *deleteIdentityProviderTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	d.Admin, err = d.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, d.GrantRoleAfterFunc, err = d.GrantRoleForUserWithAfter(d.Admin.ID.String(), adminCredentials)
	d.apiTest.SetHeader("Authorization", "Bearer "+d.AccessToken)

	return err
}

func (d *deleteIdentityProviderTest) thereIsIdentityProviderWithTheFollowingDetails(idpTable *godog.Table) error {
	idpJSON, err := d.apiTest.ReadRow(idpTable, nil, false)
	if err != nil {
		return err
	}

	idpData := dto.IdentityProvider{}
	if err := d.apiTest.UnmarshalJSON([]byte(idpJSON), &idpData); err != nil {
		return err
	}

	d.identityProvider, err = d.DB.CreateIdentityProvider(context.Background(), db.CreateIdentityProviderParams{
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

	d.apiTest.URL += d.identityProvider.ID.String()

	return nil
}

func (d *deleteIdentityProviderTest) thereIsIdentityProviderWithId(id string) error {
	d.apiTest.URL += id
	return nil
}
func (d *deleteIdentityProviderTest) iDeleteTheIdentityProvider() error {
	d.apiTest.SendRequest()
	return nil
}

func (d *deleteIdentityProviderTest) theIdentityProviderShouldBeDeleted() error {
	if err := d.apiTest.AssertStatusCode(http.StatusNoContent); err != nil {
		return err
	}

	_, err := d.DB.GetIdentityProvider(context.Background(), d.identityProvider.ID)
	if err == nil {
		return fmt.Errorf("client is not deleted")
	}
	if !sqlcerr.Is(err, sqlcerr.ErrNoRows) {
		return err
	}

	return nil
}

func (d *deleteIdentityProviderTest) theDeleteShouldFailWithErrorMessage(message string) error {
	if err := d.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return d.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (d *deleteIdentityProviderTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		d.apiTest.URL = "/v1/identityProviders/"
		d.apiTest.Method = http.MethodDelete
		d.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = d.DB.DeleteUser(ctx, d.Admin.ID)
		_, _ = d.DB.DeleteIdentityProvider(ctx, d.identityProvider.ID)
		_ = d.GrantRoleAfterFunc()
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, d.iAmLoggedInAsAdminUser)
	ctx.Step(`^I delete the identity provider$`, d.iDeleteTheIdentityProvider)
	ctx.Step(`^The delete should fail with error message "([^"]*)"$`, d.theDeleteShouldFailWithErrorMessage)
	ctx.Step(`^The identity provider should be deleted$`, d.theIdentityProviderShouldBeDeleted)
	ctx.Step(`^There is identity provider with the following details$`, d.thereIsIdentityProviderWithTheFollowingDetails)
	ctx.Step(`^There is identity provider with id "([^"]*)"$`, d.thereIsIdentityProviderWithId)
}
