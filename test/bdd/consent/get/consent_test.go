package consent

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sso/platform/utils"
	"sso/test"
	"testing"

	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type getConsentTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	redisSeeder seed.RedisDB
	redisModel  seed.RedisModel
	client      dto.Client
	consent     dto.Consent
	User        db.User
	scopes      []db.Scope
}

func TestGetConsentByID(t *testing.T) {
	a := &getConsentTest{}
	a.TestInstance = test.Initiate("../../../../")
	a.redisSeeder = seed.RedisDB{DB: a.Redis}
	a.apiTest.InitializeTest(t, "Get consent by id test", "features/consent.feature", a.InitializeScenario)
}

func (g *getConsentTest) iAmLoggedInWithCredentials(credentials *godog.Table) error {
	var err error
	g.User, err = g.Authenticate(credentials)
	if err != nil {
		return err
	}
	return nil
}

func (g *getConsentTest) thereAreRegisteredScopesWithTheFollowingDetails(scopes *godog.Table) error {
	scopesData, err := g.apiTest.ReadRows(scopes, nil, false)
	if err != nil {
		return err
	}
	var scopesStruct []dto.Scope
	if err := g.apiTest.UnmarshalJSONAt([]byte(scopesData), "", &scopesStruct); err != nil {
		return err
	}
	for _, scope := range scopesStruct {
		savedScope, err := g.DB.CreateScope(context.Background(), db.CreateScopeParams{
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
		g.scopes = append(g.scopes, savedScope)
	}
	return nil
}

func (g *getConsentTest) thereIsAClientWithTheFollowingDetails(client *godog.Table) error {
	body, err := g.apiTest.ReadRow(client, []src.Type{
		{
			Column: "redirect_uris",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	if err := g.apiTest.UnmarshalJSONAt([]byte(body), "", &g.client); err != nil {
		return err
	}

	clientData, err := g.DB.CreateClient(context.Background(), db.CreateClientParams{
		Name:         g.client.Name,
		RedirectUris: utils.ArrayToString(g.client.RedirectURIs),
		Secret:       g.client.Secret,
		Scopes:       g.client.Scopes,
		ClientType:   g.client.ClientType,
		LogoUrl:      g.client.LogoURL,
	})
	if err != nil {
		return err
	}
	g.client.ID = clientData.ID
	return nil
}

func (g *getConsentTest) iHaveAConsentWithTheFollowingDetails(consent *godog.Table) error {
	consentData, err := g.apiTest.ReadRow(consent, []src.Type{
		{
			Column: "approved",
			Kind:   src.Bool,
		},
	}, false)
	if err != nil {
		return err
	}
	err = g.apiTest.UnmarshalJSONAt([]byte(consentData), "", &g.consent)
	if err != nil {
		return err
	}
	g.consent.ClientID = g.client.ID
	consentValue, err := json.Marshal(g.consent)
	if err != nil {
		return err
	}
	g.redisModel = seed.RedisModel{
		Key:   "consent:" + g.consent.ID.String(),
		Value: string(consentValue),
	}
	err = g.redisSeeder.Feed(g.redisModel)
	if err != nil {
		return err
	}
	return nil
}

func (g *getConsentTest) iHaveAConsentWithID(consentID string) error {
	g.apiTest.URL += "/" + consentID
	return nil
}

func (g *getConsentTest) iRequestConsentData() error {
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	g.apiTest.SendRequest()
	return nil
}

func (g *getConsentTest) iShouldGetValidConsentData() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	var consentResponse dto.ConsentResponse
	err := g.apiTest.UnmarshalResponseBodyPath("data", &consentResponse)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(consentResponse.Approved, g.consent.Approved); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(consentResponse.ClientName, g.client.Name); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(consentResponse.ClientType, g.client.ClientType); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(consentResponse.ClientTrusted, false); err != nil { // FIXME: should be actually implemented
		return err
	}
	if err := g.apiTest.AssertEqual(consentResponse.ClientID, g.client.ID); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(consentResponse.UserID, g.User.ID); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(consentResponse.ClientLogo, g.client.LogoURL); err != nil {
		return err
	}
	var scopes []string
	for _, scope := range consentResponse.Scopes {
		scopes = append(scopes, scope.Name)
	}
	if err := g.apiTest.AssertEqual(utils.ArrayToString(scopes), g.consent.Scope); err != nil {
		return err
	}
	return nil
}
func (g *getConsentTest) iShouldGetErrorWithMessageAndFieldError(message, fieldError string) error {
	if message != "" {
		if err := g.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
			return err
		}
	}
	if fieldError != "" {
		if err := g.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", fieldError); err != nil {
			return err
		}
	}
	return nil
}

func (g *getConsentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.URL = "/v1/oauth/consent"
		g.apiTest.Method = "GET"
		g.apiTest.SetHeader("Content-Type", "application/json")
		g.apiTest.InitializeServer(g.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.User.ID)
		_, _ = g.DB.DeleteClient(ctx, g.client.ID)
		for _, scope := range g.scopes {
			_, _ = g.DB.DeleteScope(ctx, scope.Name)
		}
		_ = g.redisSeeder.Starve(g.redisModel)
		return ctx, nil
	})

	ctx.Step(`^I am logged in with credentials$`, g.iAmLoggedInWithCredentials)
	ctx.Step(`^There are registered scopes with the following details$`, g.thereAreRegisteredScopesWithTheFollowingDetails)
	ctx.Step(`^There is a client with the following details$`, g.thereIsAClientWithTheFollowingDetails)
	ctx.Step(`^I have a consent with ID "([^"]*)"$`, g.iHaveAConsentWithID)
	ctx.Step(`^I request consent Data$`, g.iRequestConsentData)
	ctx.Step(`^I should get valid consent data$`, g.iShouldGetValidConsentData)
	ctx.Step(`^I have a consent with the following details$`, g.iHaveAConsentWithTheFollowingDetails)
	ctx.Step(`^I should get error with message "([^"]*)" and field error "([^"]*)"$`, g.iShouldGetErrorWithMessageAndFieldError)
}
