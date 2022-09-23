package get_authorized_clients

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/test"
	"testing"
	"time"
)

type GetAuthorizedClientsTest struct {
	test.TestInstance
	apiTest           src.ApiTest
	user              db.User
	clients           []db.Client
	filters           []request_models.Filter
	authRefreshTokens []db.RefreshToken
}

func TestGetAuthorizedClients(t *testing.T) {
	c := GetAuthorizedClientsTest{}
	c.apiTest.URL = "/v1/oauth/getAuthorizedClients"
	c.apiTest.Method = http.MethodGet
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "get authorized clients test", "features/get_authorized_clients.feature", c.InitializeScenario)
}

func (g *GetAuthorizedClientsTest) iAmLoggedInAsTheFollowingUser(userCredentials *godog.Table) error {
	var err error
	g.user, err = g.Authenticate(userCredentials)
	if err != nil {
		return err
	}
	return nil
}

func (g *GetAuthorizedClientsTest) iHaveGivenAuthorizationForTheFollowingClients(clients *godog.Table) error {
	// register clients
	authorizationJSON, err := g.apiTest.ReadRows(clients, []src.Type{
		{
			WithName: "client",
			Columns:  []string{"name", "client_type", "redirect_urls", "scopes", "logo_url"},
			Kind:     src.Object,
		},
	}, true)
	if err != nil {
		return err
	}
	var authorizationData []struct {
		Client        db.CreateClientParams
		GrantedScopes string
	}
	err = g.apiTest.UnmarshalJSON([]byte(authorizationJSON), &authorizationData)
	if err != nil {
		return err
	}
	for _, v := range authorizationData {
		client, err := g.DB.CreateClient(context.Background(), v.Client)
		if err != nil {
			return err
		}
		refreshToken, err := g.DB.SaveRefreshToken(context.Background(), db.SaveRefreshTokenParams{
			ExpiresAt: time.Now().Add(5 * time.Minute),
			UserID:    g.user.ID,
			Scope: sql.NullString{
				String: v.GrantedScopes,
				Valid:  true,
			},
			RedirectUri: sql.NullString{
				String: v.Client.RedirectUris,
				Valid:  true,
			},
			ClientID:     client.ID,
			RefreshToken: "some_refresh_token",
			Code:         "some_code",
		})
		if err != nil {
			return err
		}
		g.clients = append(g.clients, client)
		g.authRefreshTokens = append(g.authRefreshTokens, refreshToken)
	}

	return nil
}

func (g *GetAuthorizedClientsTest) iRequestToGetAuthorizedClientsWithTheFollowingFilter(filter *godog.Table) error {
	filterJSON, err := g.apiTest.ReadRows(filter, nil, false)
	if err != nil {
		return err
	}
	var filterData request_models.Filter
	err = g.apiTest.UnmarshalJSONAt([]byte(filterJSON), "0", &filterData)
	if err != nil {
		return err
	}
	g.filters = append(g.filters, filterData)

	g.apiTest.SetQueryParam("filter", filterJSON)
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	g.apiTest.SendRequest()
	return nil
}

func (g *GetAuthorizedClientsTest) iShouldGetErrorMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := g.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (g *GetAuthorizedClientsTest) iShouldGetTheListOfAuthorizedClientsThatPassMyFilter() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var authClientsResponse []dto.AuthorizedClientsResponse
	err := g.apiTest.UnmarshalResponseBodyPath("data", authClientsResponse)
	if err != nil {
		return err
	}

	for _, client := range authClientsResponse {
		found := false
		for _, v := range g.clients {
			if client.ID.String() == v.ID.String() {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected client: %v", client)
		}
	}
	return nil
}

func (g *GetAuthorizedClientsTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for _, v := range g.clients {
			_, _ = g.DB.DeleteClient(ctx, v.ID)
		}
		for _, v := range g.authRefreshTokens {
			_ = g.DB.RemoveRefreshToken(ctx, v.RefreshToken)
		}
		_, _ = g.DB.DeleteUser(ctx, g.user.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in as the following user$`, g.iAmLoggedInAsTheFollowingUser)
	ctx.Step(`^I have given authorization for the following clients$`, g.iHaveGivenAuthorizationForTheFollowingClients)
	ctx.Step(`^I request to get authorized clients with the following filter$`, g.iRequestToGetAuthorizedClientsWithTheFollowingFilter)
	ctx.Step(`^I should get error message "([^"]*)"$`, g.iShouldGetErrorMessage)
	ctx.Step(`^I should get the list of authorized clients that pass my filter$`, g.iShouldGetTheListOfAuthorizedClientsThatPassMyFilter)
}
