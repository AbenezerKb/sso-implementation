package get_identity_provider

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

type getIdentityProviderTest struct {
	test.TestInstance
	apiTest          src.ApiTest
	identityProvider db.IdentityProvider
	Admin            db.User
}

func TestGetIdentityProvider(t *testing.T) {
	g := &getIdentityProviderTest{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest = src.ApiTest{
		Server: g.Server,
	}

	g.apiTest.InitializeTest(t, "get identity provider test", "features/get_identity_provider.feature", g.InitializeScenario)
}

func (g *getIdentityProviderTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	g.Admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, g.GrantRoleAfterFunc, err = g.GrantRoleForUserWithAfter(g.Admin.ID.String(), adminCredentials)
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)

	return err
}

func (g *getIdentityProviderTest) thereIsIdentityProviderWithTheFollowingDetails(idpTable *godog.Table) error {
	idpJSON, err := g.apiTest.ReadRow(idpTable, nil, false)
	if err != nil {
		return err
	}

	idpData := dto.IdentityProvider{}
	if err := g.apiTest.UnmarshalJSON([]byte(idpJSON), &idpData); err != nil {
		return err
	}

	g.identityProvider, err = g.DB.CreateIdentityProvider(context.Background(), db.CreateIdentityProviderParams{
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

	g.apiTest.URL += g.identityProvider.ID.String()

	return nil
}

func (g *getIdentityProviderTest) iHaveIdentityProviderWithId(idPid string) error {
	g.apiTest.URL += idPid

	return nil
}

func (g *getIdentityProviderTest) iGetTheIdentityProvider() error {
	g.apiTest.SendRequest()
	return nil
}

func (g *getIdentityProviderTest) iShouldSuccessfullyGetTheIdentityProvider() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	fetchedIdP := dto.IdentityProvider{}
	err := g.apiTest.UnmarshalResponseBodyPath("data", &fetchedIdP)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.ClientID, g.identityProvider.ClientID); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.ClientSecret, g.identityProvider.ClientSecret); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.LogoURI, g.identityProvider.LogoUrl.String); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.TokenEndpointURI, g.identityProvider.TokenEndpointUrl); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.UserInfoEndpointURI, g.identityProvider.UserInfoEndpointUrl.String); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.Name, g.identityProvider.Name); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.AuthorizationURI, g.identityProvider.AuthorizationUri); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedIdP.RedirectURI, g.identityProvider.RedirectUri); err != nil {
		return err
	}

	return nil
}

func (g *getIdentityProviderTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return g.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (g *getIdentityProviderTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.URL = "/v1/identityProviders/"
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.Admin.ID)
		_, _ = g.DB.DeleteIdentityProvider(ctx, g.identityProvider.ID)
		_ = g.GrantRoleAfterFunc()
		return ctx, nil
	})
	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I Get the identity provider$`, g.iGetTheIdentityProvider)
	ctx.Step(`^I have identity provider with id "([^"]*)"$`, g.iHaveIdentityProviderWithId)
	ctx.Step(`^I should successfully get the identity provider$`, g.iShouldSuccessfullyGetTheIdentityProvider)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, g.thenIShouldGetErrorWithMessage)
	ctx.Step(`^There is identity provider with the following details$`, g.thereIsIdentityProviderWithTheFollowingDetails)
}
