package codegrantflow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/state"
	"sso/platform/utils"
	"sso/test"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type issueAccessTokenCodeGrantTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	AccessToken dto.TokenResponse
	client      db.Client
	user        db.User
	redisSeeder seed.RedisDB
	authCode    seed.RedisModel
}

func TestIssueAccessTokenCodeGrant(t *testing.T) {

	i := &issueAccessTokenCodeGrantTest{}

	i.TestInstance = test.Initiate("../../../../../")
	i.redisSeeder = seed.RedisDB{
		DB: i.Redis,
	}
	i.apiTest.InitializeServer(i.Server)
	i.apiTest.InitializeTest(t, "issue acces token", "features/code_grant_flow.feature", i.InitializeScenario)
}
func (i *issueAccessTokenCodeGrantTest) theUserGrantedAccessToTheClient(arg1 *godog.Table) error {
	code, err := i.apiTest.ReadCellString(arg1, "code")
	if err != nil {
		return err
	}

	authCode := dto.AuthCode{
		Code:        code,
		Scope:       "openid",
		RedirectURI: "https://www.google.com",
		ClientID:    i.client.ID,
		UserID:      i.user.ID,
	}

	authCodeValue, err := json.Marshal(authCode)
	if err != nil {
		return err
	}
	i.authCode = seed.RedisModel{
		Key:      fmt.Sprintf(state.AuthCodeKey, authCode.Code),
		Value:    string(authCodeValue),
		ExpireAt: time.Duration(time.Minute * 2),
	}

	err = i.redisSeeder.Feed(i.authCode)
	if err != nil {
		return err
	}
	return nil
}

func (i *issueAccessTokenCodeGrantTest) aUserIsRegisteredOnTheSystem() error {
	var err error
	hash, err := utils.HashAndSalt(context.Background(), []byte("password"), i.Logger)
	if err != nil {
		return err
	}
	if i.user, err = i.DB.CreateUser(context.Background(), db.CreateUserParams{
		Email:      utils.StringOrNull("yonaskemon@gmail.com"),
		Password:   hash,
		FirstName:  "someone",
		MiddleName: "someone",
		LastName:   "someone",
		Phone:      "0987654321",
	}); err != nil {
		return err
	}

	return nil
}

func (i *issueAccessTokenCodeGrantTest) iHaveTheFollowingParameters(tokenParam *godog.Table) error {
	body, err := i.apiTest.ReadRow(tokenParam, nil, false)
	if err != nil {
		return err
	}
	i.apiTest.Body = body
	return nil
}

func (i *issueAccessTokenCodeGrantTest) theClientRequestForToken() error {
	i.apiTest.SetHeader("Authorization", "Basic "+basicAuth(i.client.ID.String(), i.client.Secret))
	i.apiTest.SetHeader("Content-Type", "application/json")
	i.apiTest.SendRequest()
	return nil
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func (i *issueAccessTokenCodeGrantTest) theRequestShouldFailWithFieldErrorAndMessage(fieldMessage, errorMessage string) error {
	if err := i.apiTest.AssertBodyColumn("error.message", errorMessage); err != nil {
		return err
	}
	return i.apiTest.AssertBodyColumn("error.field_error.0.description", fieldMessage)
}

func (i *issueAccessTokenCodeGrantTest) aClientIsRegisteredOnTheSystem() error {
	var err error
	if i.client, err = i.DB.CreateClient(context.Background(), db.CreateClientParams{
		RedirectUris: utils.ArrayToString([]string{"https://www.google.com"}),
		Name:         "google",
		Scopes:       "openid",
		ClientType:   "confidential",
		Secret:       utils.GenerateRandomString(25, true),
		LogoUrl:      "https://www.google.com/images/errors/robot.png",
	}); err != nil {
		return err
	}
	return nil
}

func (i *issueAccessTokenCodeGrantTest) tokenShouldSuccessfullyBeIssued() error {
	if err := i.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := i.apiTest.AssertColumnExists("data.access_token"); err != nil {
		return err
	}

	if err := i.apiTest.AssertColumnExists("data.refresh_token"); err != nil {
		return err
	}
	return nil
}

func (i *issueAccessTokenCodeGrantTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		i.apiTest.URL = "/v1/oauth/token"
		i.apiTest.Method = http.MethodPost
		i.apiTest.SetHeader("Content-Type", "application/json")

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {

		_, _ = i.DB.DeleteUser(context.Background(), i.user.ID)
		_ = i.redisSeeder.Starve(i.authCode)
		_, _ = i.Conn.Exec(ctx, "Delete from auth_histories where true")
		_, _ = i.Conn.Exec(ctx, "Delete from refresh_tokens where true")
		_, _ = i.DB.DeleteClient(context.Background(), i.client.ID)
		return ctx, nil
	})

	ctx.Step(`^I have the following parameters:$`, i.iHaveTheFollowingParameters)
	ctx.Step(`^The client request for token$`, i.theClientRequestForToken)
	ctx.Step(`^The request should fail with field error "([^"]*)" and message "([^"]*)"$`, i.theRequestShouldFailWithFieldErrorAndMessage)
	ctx.Step(`^The user granted access to the client:$`, i.theUserGrantedAccessToTheClient)
	ctx.Step(`^A client is registered on the system$`, i.aClientIsRegisteredOnTheSystem)
	ctx.Step(`^A user is registered on the system$`, i.aUserIsRegisteredOnTheSystem)
	ctx.Step(`^Token should successfully be issued$`, i.tokenShouldSuccessfullyBeIssued)
}
