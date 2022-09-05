package approve

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
	"net/http"
	"net/url"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"
)

type approveConsentTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	redisSeeder seed.RedisDB
	redisModel  seed.RedisModel
	client      dto.Client
	consent     dto.Consent
	User        db.User
	scopes      []db.Scope
}

func TestApproveConsent(t *testing.T) {
	a := &approveConsentTest{}
	a.TestInstance = test.Initiate("../../../../")
	a.redisSeeder = seed.RedisDB{DB: a.Redis}
	a.apiTest.URL = "/v1/oauth/approveConsent"
	a.apiTest.Method = "POST"
	a.apiTest.SetHeader("Content-Type", "application/json")
	a.apiTest.InitializeServer(a.Server)
	a.apiTest.InitializeTest(t, "approve consent", "features/approve_consent.feature", a.InitializeScenario)
}

func (a *approveConsentTest) iAmLoggedInWithCredentials(credentials *godog.Table) error {
	user, err := a.Authenticate(credentials)
	if err != nil {
		return err
	}
	a.User = user
	a.apiTest.SetHeader("Authorization", "Bearer "+a.AccessToken)
	return nil
}

func (a *approveConsentTest) thereAreRegisteredScopesWithTheFollowingDetails(scopes *godog.Table) error {
	scopesData, err := a.apiTest.ReadRows(scopes, nil, false)
	if err != nil {
		return err
	}
	var scopesStruct []dto.Scope
	if err := a.apiTest.UnmarshalJSONAt([]byte(scopesData), "", &scopesStruct); err != nil {
		return err
	}
	for _, scope := range scopesStruct {
		savedScope, err := a.DB.CreateScope(context.Background(), db.CreateScopeParams{
			Name:        scope.Name,
			Description: scope.Description,
			ResourceServerName: sql.NullString{
				String: scope.ResourceServerName,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
		a.scopes = append(a.scopes, savedScope)
	}
	return nil
}

func (a *approveConsentTest) thereIsAClientWithTheFollowingDetails(client *godog.Table) error {
	body, err := a.apiTest.ReadRow(client, []src.Type{
		{
			Column: "redirect_uris",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	if err := a.apiTest.UnmarshalJSONAt([]byte(body), "", &a.client); err != nil {
		return err
	}

	clientData, err := a.DB.CreateClient(context.Background(), db.CreateClientParams{
		Name:         a.client.Name,
		RedirectUris: utils.ArrayToString(a.client.RedirectURIs),
		Secret:       a.client.Secret,
		Scopes:       a.client.Scopes,
		ClientType:   a.client.ClientType,
		LogoUrl:      a.client.LogoURL,
	})
	if err != nil {
		return err
	}
	a.client.ID = clientData.ID
	return nil
}

func (a *approveConsentTest) iHaveAConsentWithTheFollowingDetails(consent *godog.Table) error {
	consentData, err := a.apiTest.ReadRow(consent, []src.Type{
		{
			Column: "approved",
			Kind:   src.Bool,
		},
	}, false)
	if err != nil {
		return err
	}
	err = a.apiTest.UnmarshalJSON([]byte(consentData), &a.consent)
	if err != nil {
		return err
	}
	a.consent.ClientID = a.client.ID
	consentValue, err := json.Marshal(a.consent)
	if err != nil {
		return err
	}
	a.redisModel = seed.RedisModel{
		Key:   "consent:" + a.consent.ID.String(),
		Value: string(consentValue),
	}
	err = a.redisSeeder.Feed(a.redisModel)
	if err != nil {
		return err
	}
	return nil
}

func (a *approveConsentTest) iRequestConsentApprovalWithId(consentID string) error {
	a.apiTest.SetBodyValue("consent_id", consentID)
	a.apiTest.AddCookie(http.Cookie{
		Name:  "opbs",
		Value: utils.GenerateNewOPBS(),
	})
	a.apiTest.SendRequest()
	return nil
}

func (a *approveConsentTest) theConsentShouldBeApproved() error {
	fmt.Println(string(a.apiTest.ResponseBody))
	if err := a.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var data dto.RedirectResponse
	if err := a.apiTest.UnmarshalResponseBodyPath("data", &data); err != nil {
		return err
	}
	redirectURL, err := url.Parse(data.Location)
	if err != nil {
		return err
	}
	queryParams := redirectURL.Query()
	code, err := a.TestInstance.CacheLayer.AuthCodeCacheLayer.GetAuthCode(context.Background(), queryParams.Get("code"))
	if err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(code.Scope, a.consent.Scope); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(code.RedirectURI, a.consent.RedirectURI); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(code.State, a.consent.State); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(code.UserID, a.User.ID); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(code.ClientID, a.client.ID); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(code.Code, queryParams.Get("code")); err != nil {
		return err
	}
	return nil
}

func (a *approveConsentTest) consentApprovalShouldFailWithMessage(message string) error {
	if err := a.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var data dto.RedirectResponse
	if err := a.apiTest.UnmarshalResponseBodyPath("data", &data); err != nil {
		return err
	}
	redirectURL, err := url.Parse(data.Location)
	if err != nil {
		return err
	}
	queryParams := redirectURL.Query()
	if err := a.apiTest.AssertEqual(queryParams.Get("error"), message); err != nil {
		return err
	}
	return nil
}

func (a *approveConsentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = a.DB.DeleteUser(ctx, a.User.ID)
		_, _ = a.DB.DeleteClient(ctx, a.client.ID)
		for _, scope := range a.scopes {
			_, _ = a.DB.DeleteScope(ctx, scope.Name)
		}
		_ = a.redisSeeder.Starve(a.redisModel)
		_ = a.Redis.FlushDB(ctx)
		return ctx, nil
	})
	ctx.Step(`^Consent approval should fail with message "([^"]*)"$`, a.consentApprovalShouldFailWithMessage)
	ctx.Step(`^I am logged in with credentials$`, a.iAmLoggedInWithCredentials)
	ctx.Step(`^I have a consent with the following details$`, a.iHaveAConsentWithTheFollowingDetails)
	ctx.Step(`^I request consent approval with id "([^"]*)"$`, a.iRequestConsentApprovalWithId)
	ctx.Step(`^The consent should be approved$`, a.theConsentShouldBeApproved)
	ctx.Step(`^There are registered scopes with the following details$`, a.thereAreRegisteredScopesWithTheFollowingDetails)
	ctx.Step(`^There is a client with the following details$`, a.thereIsAClientWithTheFollowingDetails)
}
