package refreshtoken

import (
	"context"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type refreshSSOTokenTest struct {
	test.TestInstance
	apiTest      src.ApiTest
	user         dto.User
	refreshToken dto.InternalRefreshToken
	redisSeeder  seed.RedisDB
}

func TestRefreshSSOToken(t *testing.T) {

	r := &refreshSSOTokenTest{}

	r.TestInstance = test.Initiate("../../../../")
	r.redisSeeder = seed.RedisDB{
		DB: r.Redis,
	}
	r.apiTest.InitializeServer(r.Server)
	r.apiTest.InitializeTest(t, "refresh access token", "features/refresh_token.feature", r.InitializeScenario)
}

func (r *refreshSSOTokenTest) iAmLoggedInToTheSystemAndHaveARefreshToken(internalRefreshToken *godog.Table) error {
	body, err := r.apiTest.ReadRow(internalRefreshToken, nil, false)
	if err != nil {
		return err
	}
	if err := r.apiTest.UnmarshalJSON([]byte(body), &r.refreshToken); err != nil {
		return err
	}
	rfData, err := r.DB.SaveInternalRefreshToken(context.Background(), db.SaveInternalRefreshTokenParams{
		UserID:       r.user.ID,
		ExpiresAt:    r.refreshToken.ExpiresAt,
		RefreshToken: r.refreshToken.RefreshToken,
	})
	if err != nil {
		return err
	}
	r.refreshToken.RefreshToken = rfData.RefreshToken
	return nil
}

func (r *refreshSSOTokenTest) iRefreshMyAccessTokenUsingMyRefreshToken(rfParam *godog.Table) error {
	refreshToken, err := r.apiTest.ReadCellString(rfParam, "refresh_token")
	if err != nil {
		return err
	}
	r.apiTest.AddCookie(http.Cookie{
		Name:     "ab_fen",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		MaxAge:   3600,
		HttpOnly: true,
	})

	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.SendRequest()
	return nil
}

func (r *refreshSSOTokenTest) iShouldGetANewAccessToken() error {
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

func (r *refreshSSOTokenTest) thereIsARegisteredUserOnTheSystem(user *godog.Table) error {
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

func (r *refreshSSOTokenTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/refreshToken"
		r.apiTest.Method = http.MethodGet
		r.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(context.Background(), r.user.ID)
		_ = r.DB.RemoveInternalRefreshToken(context.Background(), r.refreshToken.RefreshToken)
		return ctx, nil
	})
	ctx.Step(`^I am logged in to the system and have a refresh token:$`, r.iAmLoggedInToTheSystemAndHaveARefreshToken)
	ctx.Step(`^I refresh my access token using my refresh token$`, r.iRefreshMyAccessTokenUsingMyRefreshToken)
	ctx.Step(`^I should get a new access token$`, r.iShouldGetANewAccessToken)
	ctx.Step(`^There is a registered user on the system:$`, r.thereIsARegisteredUserOnTheSystem)
}
