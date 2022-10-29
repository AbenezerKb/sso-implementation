package identityproviders

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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

	u.apiTest.InitializeTest(t, "get identity provider's test", "features/identity_providers.feature", u.InitializeScenario)
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

func (g *getIdentityProvidersTest) iRequestToGetAllTheIdentityProviders() error {
	g.apiTest.SendRequest()
	return nil
}

func (g *getIdentityProvidersTest) iShouldGetAllTheIdentityProviders() error {
	var responseIdPs []dto.IdentityProvider

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	err := g.apiTest.UnmarshalResponseBodyPath("data", &responseIdPs)
	if err != nil {
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
		g.apiTest.URL = "/v1/registeredIdentityProviders"
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for i := 0; i < len(g.identityProviders); i++ {
			_, _ = g.DB.DeleteIdentityProvider(ctx, g.identityProviders[i].ID)
		}

		return ctx, nil
	})

	ctx.Step(`^I request to get all the identity providers$`, g.iRequestToGetAllTheIdentityProviders)
	ctx.Step(`^I should get all the identity providers$`, g.iShouldGetAllTheIdentityProviders)
	ctx.Step(`^There are identity provider with the following details$`, g.thereAreIdentityProviderWithTheFollowingDetails)
}
