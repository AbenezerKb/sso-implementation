package approve

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/cucumber/godog"
	"github.com/joomcode/errorx"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
	"net/http"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"
)

type rejectConsentTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	redisSeeder seed.RedisDB
	redisModel  seed.RedisModel
	client      dto.Client
	consent     dto.Consent
	User        db.User
	scopes      []db.Scope
}

func TestRejectConsent(t *testing.T) {
	a := &rejectConsentTest{}
	a.TestInstance = test.Initiate("../../../../")
	a.redisSeeder = seed.RedisDB{DB: a.Redis}
	a.apiTest.URL = "/v1/oauth/rejectConsent"
	a.apiTest.Method = "POST"
	a.apiTest.SetHeader("Content-Type", "application/json")
	a.apiTest.InitializeServer(a.Server)
	a.apiTest.InitializeTest(t, "reject consent", "features/reject_consent.feature", a.InitializeScenario)
}

func (a *rejectConsentTest) iAmLoggedInWithCredentials(credentials *godog.Table) error {
	user, err := a.Authenticate(credentials)
	if err != nil {
		return err
	}
	a.User = user
	a.apiTest.SetHeader("Authorization", "Bearer "+a.AccessToken)
	return nil
}

func (a *rejectConsentTest) thereAreRegisteredScopesWithTheFollowingDetails(scopes *godog.Table) error {
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

func (a *rejectConsentTest) thereIsAClientWithTheFollowingDetails(client *godog.Table) error {
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

func (a *rejectConsentTest) iHaveAConsentWithTheFollowingDetails(consent *godog.Table) error {
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

func (a *rejectConsentTest) iRequestConsentRejectionWithIdAndMessage(consentID string, message string) error {
	//a.apiTest.Body = `{"consent_id":"` + consentID + `"}`
	a.apiTest.SetQueryParam("consentId", consentID)
	a.apiTest.SetQueryParam("failureReason", message)
	a.apiTest.AddCookie(http.Cookie{
		Name:  "opbs",
		Value: utils.GenerateNewOPBS(),
	})
	a.apiTest.SendRequest()
	return nil
}

func (a *rejectConsentTest) theConsentShouldBeRejected() error {
	if err := a.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}
	queryParams := a.apiTest.GetRedirectURLQueryParams()
	_, err := a.TestInstance.CacheLayer.AuthCodeCacheLayer.GetAuthCode(context.Background(), queryParams["code"])
	if !errorx.IsOfType(err, errors.ErrNoRecordFound) {
		return err
	}

	if err := a.apiTest.AssertEqual(queryParams["state"], a.consent.State); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(queryParams["error"], "access_denied"); err != nil {
		return err
	}
	return nil
}

func (a *rejectConsentTest) consentRejectionShouldFailWithMessage(message string) error {
	if err := a.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}
	queryParams := a.apiTest.GetRedirectURLQueryParams()
	if err := a.apiTest.AssertEqual(queryParams["error"], message); err != nil {
		return err
	}
	return nil
}

func (a *rejectConsentTest) InitializeScenario(ctx *godog.ScenarioContext) {
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
	ctx.Step(`^Consent rejection should fail with message "([^"]*)"$`, a.consentRejectionShouldFailWithMessage)
	ctx.Step(`^I am logged in with credentials$`, a.iAmLoggedInWithCredentials)
	ctx.Step(`^I have a consent with the following details$`, a.iHaveAConsentWithTheFollowingDetails)
	ctx.Step(`^I request consent rejection with id "([^"]*)" and message "([^"]*)"$`, a.iRequestConsentRejectionWithIdAndMessage)
	ctx.Step(`^The consent should be rejected$`, a.theConsentShouldBeRejected)
	ctx.Step(`^There are registered scopes with the following details$`, a.thereAreRegisteredScopesWithTheFollowingDetails)
	ctx.Step(`^There is a client with the following details$`, a.thereIsAClientWithTheFollowingDetails)
}
