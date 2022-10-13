package get_identity_providers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getIdentityProvidersTest struct {
	test.TestInstance
	apiTest           src.ApiTest
	identityProviders []db.IdentityProvider
	Admin             db.User
	Preferences       preferenceData
}

type preferenceData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func TestGetIdentityProviders(t *testing.T) {
	u := &getIdentityProvidersTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest = src.ApiTest{
		Server: u.Server,
	}

	u.apiTest.InitializeTest(t, "get identity provider's test", "features/get_identity_providers.feature", u.InitializeScenario)
}

func (g *getIdentityProvidersTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	g.Admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, g.GrantRoleAfterFunc, err = g.GrantRoleForUserWithAfter(g.Admin.ID.String(), adminCredentials)
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)

	return err
}

func (g *getIdentityProvidersTest) thereAreIdentityProviderWithTheFollowingDetails(idPsTable *godog.Table) error {
	idPsJSON, err := g.apiTest.ReadRows(idPsTable, nil, false)
	if err != nil {
		return err
	}

	var idPsData []dto.IdentityProvider
	if err := g.apiTest.UnmarshalJSON([]byte(idPsJSON), &idPsData); err != nil {
		return err
	}

	for _, v := range idPsData {
		createdIdentityProvider, err := g.DB.CreateIdentityProvider(context.Background(), db.CreateIdentityProviderParams{
			Name:                v.Name,
			LogoUrl:             sql.NullString{String: v.LogoURI, Valid: true},
			ClientSecret:        v.ClientSecret,
			ClientID:            v.ClientID,
			RedirectUri:         v.RedirectURI,
			AuthorizationUri:    v.AuthorizationURI,
			TokenEndpointUrl:    v.TokenEndpointURI,
			UserInfoEndpointUrl: sql.NullString{String: v.UserInfoEndpointURI, Valid: true},
		})

		if err != nil {
			return err
		}
		g.identityProviders = append(g.identityProviders, createdIdentityProvider)
	}

	return nil
}

func (g *getIdentityProvidersTest) iRequestToGetAllTheIdentityProvidersWithTheFollowingPreferences(preferences *godog.Table) error {
	preferencesJSON, err := g.apiTest.ReadRow(preferences, []src.Type{
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

	err = g.apiTest.UnmarshalJSON([]byte(preferencesJSON), &g.Preferences)
	if err != nil {
		return err
	}

	g.apiTest.SetQueryParam("page", fmt.Sprintf("%d", g.Preferences.Page))
	g.apiTest.SetQueryParam("per_page", fmt.Sprintf("%d", g.Preferences.PerPage))
	g.apiTest.SendRequest()
	return nil
}

func (g *getIdentityProvidersTest) iShouldGetTheListOfIdentityProvidersThatPassMyPreferences() error {
	var responseIdPs []dto.IdentityProvider
	var metaData model.MetaData

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	err := g.apiTest.UnmarshalResponseBodyPath("meta_data", &metaData)
	if err != nil {
		return err
	}

	err = g.apiTest.UnmarshalResponseBodyPath("data", &responseIdPs)
	if err != nil {
		return err
	}

	var total int
	if g.Preferences.Page < metaData.Total/g.Preferences.PerPage {
		total = g.Preferences.PerPage
	} else if g.Preferences.Page == metaData.Total/g.Preferences.PerPage {
		total = metaData.Total % g.Preferences.PerPage
	} else {
		total = 0
	}

	if err := g.apiTest.AssertEqual(len(responseIdPs), total); err != nil {
		return err
	}

	for _, v := range responseIdPs {
		found := false
		for _, v2 := range g.identityProviders {
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

func (g *getIdentityProvidersTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.URL = "/v1/identityProviders"
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.Admin.ID)
		for i := 0; i < len(g.identityProviders); i++ {
			_, _ = g.DB.DeleteIdentityProvider(ctx, g.identityProviders[i].ID)
		}
		_ = g.GrantRoleAfterFunc()

		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get all the identity providers with the following preferences$`, g.iRequestToGetAllTheIdentityProvidersWithTheFollowingPreferences)
	ctx.Step(`^I should get the list of identity providers that pass my preferences$`, g.iShouldGetTheListOfIdentityProvidersThatPassMyPreferences)
	ctx.Step(`^There are identity provider with the following details$`, g.thereAreIdentityProviderWithTheFollowingDetails)
}
