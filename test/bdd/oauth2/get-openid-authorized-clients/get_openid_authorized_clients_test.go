package get_openid_authorized_clients

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"

	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
)

type GetOpenIDAuthorizedClientsTest struct {
	test.TestInstance
	apiTest           src.ApiTest
	user              db.User
	clients           []db.Client
	filters           []db_pgnflt.Filter
	authRefreshTokens []db.RefreshToken
}

func TestGetAuthorizedClients(t *testing.T) {
	c := GetOpenIDAuthorizedClientsTest{}
	c.apiTest.URL = "/v1/oauth/openIDAuthorizedClients"
	c.apiTest.Method = http.MethodGet
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "get openid authorized clients test", "features/get_openid_authorized_clients.feature", c.InitializeScenario)
}

func (g *GetOpenIDAuthorizedClientsTest) iAmLoggedInAsTheFollowingUser(userCredentials *godog.Table) error {
	var err error
	g.user, err = g.Authenticate(userCredentials)
	if err != nil {
		return err
	}
	return nil
}

func (g *GetOpenIDAuthorizedClientsTest) iHaveGivenAuthorizationForTheFollowingClients(clients *godog.Table) error {
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
		GrantedScopes string `json:"granted_scopes"`
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
		if strings.Contains(refreshToken.Scope.String, "openid") {
			g.clients = append(g.clients, client)
		}
		g.authRefreshTokens = append(g.authRefreshTokens, refreshToken)
	}

	return nil
}

func (g *GetOpenIDAuthorizedClientsTest) iRequestToGetOpenidAuthorizedClients() error {
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	g.apiTest.SendRequest()
	return nil
}

func (g *GetOpenIDAuthorizedClientsTest) iShouldGetTheListOfOpenidAuthorizedClients() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var authClientsResponse []dto.AuthorizedClientsResponse
	err := g.apiTest.UnmarshalResponseBodyPath("data", &authClientsResponse)
	if err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(len(authClientsResponse), len(g.clients)); err != nil {
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

func (g *GetOpenIDAuthorizedClientsTest) InitializeScenario(ctx *godog.ScenarioContext) {
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
	ctx.Step(`^I request to get openid authorized clients$`, g.iRequestToGetOpenidAuthorizedClients)
	ctx.Step(`^I should get the list of openid authorized clients$`, g.iShouldGetTheListOfOpenidAuthorizedClients)
}
