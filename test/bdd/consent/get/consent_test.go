package consent

import (
	"context"
	"encoding/json"
	"net/http"
	"sso/platform/utils"
	"sso/test"
	"testing"

	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"

	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type getConsentTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	user_id     string
	consentID   string
	redisSeeder seed.RedisDB
	redisModel  seed.RedisModel
	userData    db.User
}

func TestGetConsentByID(t *testing.T) {
	a := &getConsentTest{}
	a.TestInstance = test.Initiate("../../../../")
	a.redisSeeder = seed.RedisDB{DB: a.Redis}
	a.apiTest.InitializeTest(t, "Get consent by id test", "features/consent.feature", a.InitializeScenario)

}
func (g *getConsentTest) iAmLoggedInWithCredentials(credentials *godog.Table) error {
	return g.Authenicate(credentials)
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
	g.consentID, _ = g.apiTest.ReadCellString(consent, "consent_id")
	g.user_id, _ = g.apiTest.ReadCellString(consent, "user_id")
	scopes, _ := g.apiTest.ReadCellString(consent, "scopes")
	redirectURI, _ := g.apiTest.ReadCellString(consent, "redirect_uri")
	clientID, _ := g.apiTest.ReadCellString(consent, "client_id")
	status, _ := g.apiTest.ReadCellString(consent, "status")

	consentID, _ := uuid.Parse(g.consentID)
	userID, _ := uuid.Parse(g.user_id)
	ParsedclientID, _ := uuid.Parse(clientID)

	consents := dto.Consent{
		ID:     consentID,
		UserID: userID,
		AuthorizationRequestParam: dto.AuthorizationRequestParam{
			Scope:       scopes,
			RedirectURI: redirectURI,
			ClientID:    ParsedclientID,
		},
		Approved: bool(status == "approved"),
	}
	// marshal consents to string
	consentString, err := json.Marshal(consents)
	if err != nil {
		return err
	}
	g.redisModel = seed.RedisModel{
		Key:   "consent:" + g.consentID,
		Value: string(consentString),
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

func (g *getConsentTest) userWithID(user_id string) error {
	// seed user
	hash, err := utils.HashAndSalt(context.Background(), []byte("password"), g.Logger)
	if err != nil {
		return err
	}
	userData, err := g.DB.CreateUser(context.Background(), db.CreateUserParams{
		Phone:    "1234567890",
		Email:    utils.StringOrNull("email"),
		Password: hash,
	})
	if err != nil {
		return err
	}
	g.userData = userData
	g.apiTest.SetQueryParam("user_id", userData.ID.String())
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
		_, _ = g.DB.DeleteUser(ctx, g.userData.ID)
		_, _ = g.DB.DeleteUser(ctx, g.User.ID)
		_ = g.redisSeeder.Starve(g.redisModel)
		return ctx, nil
	})

	ctx.Step(`^I am logged in with credentials$`, g.iAmLoggedInWithCredentials)
	ctx.Step(`^I have a consent with ID "([^"]*)"$`, g.iHaveAConsentWithID)
	ctx.Step(`^I request consent Data$`, g.iRequestConsentData)
	ctx.Step(`^I should get error "([^"]*)"$`, g.iShouldGetError)
	ctx.Step(`^I should get valid consent data$`, g.iShouldGetValidConsentData)
	ctx.Step(`^Invalid user ID "([^"]*)"$`, g.invalidUserID)
	ctx.Step(`^I have a consent with the following details$`, g.iHaveAConsentWithTheFollowingDetails)
	ctx.Step(`^user with ID "([^"]*)"$`, g.userWithID)
}
