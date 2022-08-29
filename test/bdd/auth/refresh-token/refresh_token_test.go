package refreshtoken

import (
	"context"
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

type refreshSSOTokenTest struct {
	test.TestInstance
	apiTest      src.ApiTest
	user         dto.User
	refreshToken dto.InternalRefreshToken
	redisSeeder  seed.RedisDB
}

func TestIssueAccessTokenCodeGrant(t *testing.T) {

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
	if err := r.apiTest.UnmarshalJSONAt([]byte(body), "", &r.refreshToken); err != nil {
		return err
	}
	rfData, err := r.DB.SaveInternalRefreshToken(context.Background(), db.SaveInternalRefreshTokenParams{
		UserID:       r.user.ID,
		ExpiresAt:    r.refreshToken.ExpiresAt,
		Refreshtoken: r.refreshToken.Refreshtoken,
	})
	if err != nil {
		return err
	}
	r.refreshToken.Refreshtoken = rfData.Refreshtoken
	return nil
}

func (r *refreshSSOTokenTest) iRefreshMyAccessTokenUsingMyRefresh_token(rfParam *godog.Table) error {
	body, err := r.apiTest.ReadRow(rfParam, nil, false)
	if err != nil {
		return err
	}
	r.apiTest.Body = body

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

func (r *refreshSSOTokenTest) theOldRefresh_tokenShouldBeDeleted() error {
	_, err := r.DB.GetInternalRefreshToken(context.TODO(), r.refreshToken.Refreshtoken)

	if !sqlcerr.Is(err, sqlcerr.ErrNoRows) {
		return errors.New("old refresh token not deleted.")
	}
	return nil
}

func (r *refreshSSOTokenTest) thereIsARegeisteredUserOnTheSystem(user *godog.Table) error {
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
		r.apiTest.URL = "/v1/refreshtoken"
		r.apiTest.Method = http.MethodPost
		r.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(context.Background(), r.user.ID)
		_ = r.DB.RemoveInternalRefreshToken(context.Background(), r.refreshToken.Refreshtoken)
		return ctx, nil
	})
	ctx.Step(`^I am logged in to the system and have a refresh token:$`, r.iAmLoggedInToTheSystemAndHaveARefreshToken)
	ctx.Step(`^I refresh my access token using my refresh_token$`, r.iRefreshMyAccessTokenUsingMyRefresh_token)
	ctx.Step(`^I should get a new access token$`, r.iShouldGetANewAccessToken)
	ctx.Step(`^The old refresh_token should be deleted$`, r.theOldRefresh_tokenShouldBeDeleted)
	ctx.Step(`^There is a regeistered user on the system:$`, r.thereIsARegeisteredUserOnTheSystem)
}
