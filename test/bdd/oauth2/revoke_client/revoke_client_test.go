package revoke_client

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"
	"time"
)

type revokeClientTest struct {
	test.TestInstance
	apiTest      src.ApiTest
	client       db.Client
	refreshToken db.RefreshToken
	user         db.User
}

func TestRevokeClient(t *testing.T) {

	r := &revokeClientTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.URL = "/v1/oauth/revokeClient"
	r.apiTest.Method = "POST"
	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "revoke client", "features/revoke_client.feature", r.InitializeScenario)
}

func (r *revokeClientTest) iAmLoggedInWithTheFollowingCredentials(credentials *godog.Table) error {
	user, err := r.Authenticate(credentials)
	if err != nil {
		return err
	}
	r.user = user
	return nil
}

func (r *revokeClientTest) iHaveGivenAccessToTheFollowingClient(clientData *godog.Table) error {
	clientDataJSON, err := r.apiTest.ReadRow(clientData, []src.Type{
		{
			Column: "redirect_uris",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	var clientDataDTO dto.Client
	err = r.apiTest.UnmarshalJSON([]byte(clientDataJSON), &clientDataDTO)
	if err != nil {
		return err
	}

	fmt.Println("faben", clientDataDTO)
	// register client
	r.client, err = r.DB.CreateClient(context.Background(), db.CreateClientParams{
		Name:         clientDataDTO.Name,
		ClientType:   clientDataDTO.ClientType,
		RedirectUris: utils.ArrayToString(clientDataDTO.RedirectURIs),
		Scopes:       clientDataDTO.Scopes,
		Secret:       clientDataDTO.Secret,
		LogoUrl:      clientDataDTO.LogoURL,
	})
	if err != nil {
		return err
	}
	// create refresh token on behalf of the client
	r.refreshToken, err = r.DB.SaveRefreshToken(context.Background(), db.SaveRefreshTokenParams{
		ExpiresAt: time.Now().Add(10 * time.Minute),
		UserID:    r.user.ID,
		Scope: sql.NullString{
			String: r.client.Scopes,
			Valid:  true,
		},
		RedirectUri: sql.NullString{
			String: utils.StringToArray(r.client.RedirectUris)[0],
			Valid:  true,
		},
		ClientID:     r.client.ID,
		RefreshToken: utils.GenerateRandomString(10, false),
		Code:         utils.GenerateRandomString(10, false),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *revokeClientTest) iRequestToRevokeAccessToTheClient() error {
	r.apiTest.SetHeader("Authorization", "Bearer "+r.AccessToken)
	r.apiTest.SetBodyValue("client_id", r.client.ID)
	r.apiTest.SendRequest()
	return nil
}

func (r *revokeClientTest) theClientShouldNoLongerHaveAccessToMyData() error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	_, err := r.DB.GetRefreshToken(context.Background(), r.refreshToken.RefreshToken)
	if !sqlcerr.Is(err, sqlcerr.ErrNoRows) {
		return fmt.Errorf("got %v, expected %v", err, sqlcerr.ErrNoRows)
	}
	return nil
}
func (r *revokeClientTest) myActionShouldBeRecorded() error {
	record, err := r.DB.GetLastAuthHistory(context.Background(), db.GetLastAuthHistoryParams{
		UserID:   r.user.ID,
		ClientID: r.client.ID,
	})
	if err != nil {
		return err
	}
	// TODO: may be check the time the history was recorded
	if err := r.apiTest.AssertEqual(record.Status, constant.Revoke); err != nil {
		return err
	}
	return nil
}
func (r *revokeClientTest) iRequestToRevokeAccessToTheClientWithId(clientID string) error {
	r.apiTest.SetHeader("Authorization", "Bearer "+r.AccessToken)
	r.apiTest.SetBodyValue("client_id", clientID)
	r.apiTest.SendRequest()
	return nil
}

func (r *revokeClientTest) myRequestFailsWithFieldError(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (r *revokeClientTest) myRequestFailsWithErrorMessage(message string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (r *revokeClientTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(ctx, r.user.ID)
		_, _ = r.DB.DeleteClient(ctx, r.client.ID)
		_ = r.DB.RemoveRefreshToken(ctx, r.refreshToken.RefreshToken)
		_, _ = r.Conn.Exec(ctx, "delete from auth_history where true")
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, r.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I have given access to the following client$`, r.iHaveGivenAccessToTheFollowingClient)
	ctx.Step(`^I request to revoke access to the client$`, r.iRequestToRevokeAccessToTheClient)
	ctx.Step(`^I request to revoke access to the client with id "([^"]*)"$`, r.iRequestToRevokeAccessToTheClientWithId)
	ctx.Step(`^The client should no longer have access to my data$`, r.theClientShouldNoLongerHaveAccessToMyData)
	ctx.Step(`^My action should be recorded$`, r.myActionShouldBeRecorded)
	ctx.Step(`^My request fails with field error "([^"]*)"$`, r.myRequestFailsWithFieldError)
	ctx.Step(`^My request fails with error message "([^"]*)"$`, r.myRequestFailsWithErrorMessage)
}
