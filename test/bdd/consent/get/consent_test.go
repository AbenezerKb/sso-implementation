package consent

import (
	"context"
	"encoding/json"
	"fmt"
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
func (g *getConsentTest) iHaveAConsentWithID(consentID string) error {
	g.apiTest.URL += "/" + consentID
	return nil
}

func (g *getConsentTest) iRequestConsentData() error {
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	g.apiTest.SendRequest()
	return nil
}

func (g *getConsentTest) iShouldGetError(errMsg string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	if err := g.apiTest.AssertStringValueOnPathInResponse("error.message", errMsg); err != nil {
		return err
	}

	return nil
}

func (g *getConsentTest) iShouldGetValidConsentData() error {
	fmt.Println(string(g.apiTest.ResponseBody))
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := g.apiTest.AssertColumnExists("client_id"); err != nil {
		return err
	}
	if err := g.apiTest.AssertColumnExists("scopes"); err != nil {
		return err
	}
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
	fmt.Println(g.consent.ID)
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

func (g *getConsentTest) invalidUserID(user_id string) error {
	g.apiTest.SetQueryParam("user_id", user_id)
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
		_ = g.redisSeeder.Starve(g.redisModel)
		return ctx, nil
	})

	ctx.Step(`^I am logged in with credentials$`, g.iAmLoggedInWithCredentials)
	ctx.Step(`^There is a client with the following details$`, g.thereIsAClientWithTheFollowingDetails)
	ctx.Step(`^I have a consent with ID "([^"]*)"$`, g.iHaveAConsentWithID)
	ctx.Step(`^I request consent Data$`, g.iRequestConsentData)
	ctx.Step(`^I should get error "([^"]*)"$`, g.iShouldGetError)
	ctx.Step(`^I should get valid consent data$`, g.iShouldGetValidConsentData)
	ctx.Step(`^Invalid user ID "([^"]*)"$`, g.invalidUserID)
	ctx.Step(`^I have a consent with the following details$`, g.iHaveAConsentWithTheFollowingDetails)
}
