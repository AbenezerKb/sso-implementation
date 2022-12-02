package refreshtoken

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type refreshClientTokenTest struct {
	test.TestInstance
	apiTest             src.ApiTest
	AccessToken         dto.TokenResponse
	client              dto.Client
	user                dto.User
	refreshToken        dto.RefreshToken
	expiredRefreshToken dto.RefreshToken
	redisSeeder         seed.RedisDB
}

func TestRefreshAccessToken(t *testing.T) {

	r := &refreshClientTokenTest{}

	r.TestInstance = test.Initiate("../../../../../")
	r.redisSeeder = seed.RedisDB{
		DB: r.Redis,
	}
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "refresh access token", "features/refresh_token.feature", r.InitializeScenario)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func (r *refreshClientTokenTest) theRequestShouldFailWithFieldError(msg string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", msg); err != nil {
		return err
	}
	return nil
}
func (r *refreshClientTokenTest) theRequestShouldFailWithErrorMessage(msg string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusUnauthorized); err != nil {
		return err
	}
	if err := r.apiTest.AssertStringValueOnPathInResponse("error.message", msg); err != nil {
		return err
	}
	return nil
}
func (r *refreshClientTokenTest) thereIsARegeisteredUserOnTheSystem(user *godog.Table) error {
	body, err := r.apiTest.ReadRow(user, nil, false)
	if err != nil {
		return err
	}
	if err := r.apiTest.UnmarshalJSONAt([]byte(body), "", &r.user); err != nil {
		return err
	}
	hash, err := utils.HashAndSalt(context.Background(), []byte(r.user.Password), r.Logger)
	if err != nil {
		return err
	}
	userData, err := r.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:  r.user.FirstName,
		MiddleName: r.user.MiddleName,
		LastName:   r.user.LastName,
		Password:   hash,
		Email:      utils.StringOrNull(r.user.Email),
		Phone:      r.user.Phone,
	})
	if err != nil {
		return err
	}
	r.user.ID = userData.ID
	return nil
}

func (r *refreshClientTokenTest) thereIsAClientOnTheSystem(client *godog.Table) error {
	body, err := r.apiTest.ReadRow(client, []src.Type{
		{
			Column: "redirect_uris",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	if err := r.apiTest.UnmarshalJSONAt([]byte(body), "", &r.client); err != nil {
		return err
	}

	clientData, err := r.DB.CreateClient(context.Background(), db.CreateClientParams{
		Name:         r.client.Name,
		RedirectUris: utils.ArrayToString(r.client.RedirectURIs),
		Secret:       r.client.Secret,
		Scopes:       r.client.Scopes,
		ClientType:   r.client.ClientType,
		LogoUrl:      r.client.LogoURL,
	})
	if err != nil {
		return err
	}
	r.client.ID = clientData.ID
	return nil

}
func (r *refreshClientTokenTest) iHaveAnExpiredRefreshToken(rfParam *godog.Table) error {
	body, err := r.apiTest.ReadRow(rfParam, nil, false)
	if err != nil {
		return err
	}
	if err := r.apiTest.UnmarshalJSONAt([]byte(body), "", &r.expiredRefreshToken); err != nil {
		return err
	}

	rfData, err := r.DB.SaveRefreshToken(context.Background(), db.SaveRefreshTokenParams{
		UserID:       r.user.ID,
		ClientID:     r.client.ID,
		Scope:        utils.StringOrNull(r.expiredRefreshToken.Scope),
		RefreshToken: r.expiredRefreshToken.RefreshToken,
		RedirectUri:  utils.StringOrNull(utils.ArrayToString(r.client.RedirectURIs)),
		ExpiresAt:    r.expiredRefreshToken.ExpiresAt,
	})
	if err != nil {
		return err
	}
	r.expiredRefreshToken = dto.RefreshToken{
		ID:           rfData.ID,
		UserID:       rfData.UserID,
		ClientID:     rfData.ClientID,
		Scope:        rfData.Scope.String,
		RedirectUri:  rfData.RedirectUri.String,
		ExpiresAt:    rfData.ExpiresAt,
		RefreshToken: rfData.RefreshToken,
	}
	return nil
}

func (r *refreshClientTokenTest) iRefreshTheAccessToken(rfParam *godog.Table) error {
	body, err := r.apiTest.ReadRow(rfParam, nil, false)
	if err != nil {
		return err
	}
	r.apiTest.Body = body

	r.apiTest.SetHeader("Authorization", "Basic "+basicAuth(r.client.ID.String(), r.client.Secret))
	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.SendRequest()
	return nil
}
func (r *refreshClientTokenTest) theUserGrantsAccessToTheClient(rfParam *godog.Table) error {
	body, err := r.apiTest.ReadRow(rfParam, nil, false)
	if err != nil {
		return err
	}
	if err := r.apiTest.UnmarshalJSONAt([]byte(body), "", &r.refreshToken); err != nil {
		return err
	}

	rfData, err := r.DB.SaveRefreshToken(context.Background(), db.SaveRefreshTokenParams{
		UserID:       r.user.ID,
		ClientID:     r.client.ID,
		Scope:        utils.StringOrNull(r.refreshToken.Scope),
		RefreshToken: r.refreshToken.RefreshToken,
		RedirectUri:  utils.StringOrNull(utils.ArrayToString(r.client.RedirectURIs)),
		ExpiresAt:    r.refreshToken.ExpiresAt,
	})
	if err != nil {
		return err
	}
	r.refreshToken = dto.RefreshToken{
		ID:           rfData.ID,
		UserID:       rfData.UserID,
		ClientID:     rfData.ClientID,
		Scope:        rfData.Scope.String,
		RedirectUri:  rfData.RedirectUri.String,
		ExpiresAt:    rfData.ExpiresAt,
		RefreshToken: rfData.RefreshToken,
	}
	return nil
}
func (r *refreshClientTokenTest) iShouldGetANewAccessTokenWithANewRefreshToken() error {
	if err := r.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	if err := r.apiTest.AssertColumnExists("data.access_token"); err != nil {
		return err
	}

	if err := r.apiTest.AssertColumnExists("data.refresh_token"); err != nil {
		return err
	}
	if err := r.apiTest.AssertColumnExists("data.expires_in"); err != nil {
		return err
	}
	if err := r.apiTest.AssertColumnExists("data.token_type"); err != nil {
		return err
	}
	return nil
}

func (r *refreshClientTokenTest) theOldRefreshTokenShouldBeDeleted() error {
	_, err := r.DB.GetRefreshToken(context.TODO(), r.refreshToken.RefreshToken)

	if !sqlcerr.Is(err, sqlcerr.ErrNoRows) {
		return errors.New("old refresh token not deleted.")
	}
	return nil
}
func (r *refreshClientTokenTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/oauth/token"
		r.apiTest.Method = http.MethodPost
		r.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteClient(context.Background(), r.client.ID)
		_, _ = r.DB.DeleteUser(context.Background(), r.user.ID)
		_ = r.DB.RemoveRefreshToken(context.Background(), r.AccessToken.RefreshToken)
		return ctx, nil
	})

	ctx.Step(`^I have an expired refresh token:$`, r.iHaveAnExpiredRefreshToken)
	ctx.Step(`^I refresh the access token:$`, r.iRefreshTheAccessToken)
	ctx.Step(`^I should get a new access token with the old refresh token$`, r.iShouldGetANewAccessTokenWithANewRefreshToken)
	//ctx.Step(`^The old refresh token should be deleted$`, r.theOldRefreshTokenShouldBeDeleted)
	ctx.Step(`^The request should fail with error message "([^"]*)":$`, r.theRequestShouldFailWithErrorMessage)
	ctx.Step(`^The request should fail with field error "([^"]*)":$`, r.theRequestShouldFailWithFieldError)
	ctx.Step(`^The user grants access to the client:$`, r.theUserGrantsAccessToTheClient)
	ctx.Step(`^There is a client on the system:$`, r.thereIsAClientOnTheSystem)
	ctx.Step(`^There is a regeistered user on the system:$`, r.thereIsARegeisteredUserOnTheSystem)

}
