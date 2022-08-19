package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"

	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type loginTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	user        *dto.User
	redisSeeder seed.RedisDB
}

func TestLogin(t *testing.T) {
	a := &loginTest{}
	a.TestInstance = test.Initiate("../../../")
	// create redis seeder
	a.redisSeeder = seed.RedisDB{DB: a.Redis}
	a.apiTest.InitializeTest(t, "Login test", "features/login.feature", a.InitializeScenario)

}

func (l *loginTest) iAmARegisteredUserWithDetails(userTable *godog.Table) error {
	user, err := l.apiTest.ReadRow(userTable, nil, false)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(user), &l.user)
	if err != nil {
		return err
	}
	hash, err := utils.HashAndSalt(context.Background(), []byte(l.user.Password), l.Logger)
	if err != nil {
		return err
	}
	userData, err := l.DB.CreateUser(context.Background(), db.CreateUserParams{
		Phone:    l.user.Phone,
		Email:    utils.StringOrNull(l.user.Email),
		Password: hash,
	})
	if err != nil {
		return err
	}
	l.user.ID = userData.ID
	return nil
}

func (l *loginTest) iFillTheFollowingDetails(loginInfo *godog.Table) error {
	// set otp to redis
	phone, err := l.apiTest.ReadCell(loginInfo, "phone", nil)
	if err != nil {
		return err
	}
	otp, err := l.apiTest.ReadCell(loginInfo, "otp", nil)
	if err != nil {
		return err
	}
	err = l.redisSeeder.Feed(seed.RedisModel{
		Key:   fmt.Sprintf("%s", phone),
		Value: fmt.Sprintf("%s", otp),
	})
	if err != nil {
		return err
	}
	body, err := l.apiTest.ReadRow(loginInfo, nil, false)
	if err != nil {
		return err
	}
	l.apiTest.Body = body
	return nil
}

func (l *loginTest) iSubmitTheRegistrationForm() error {
	l.apiTest.SendRequest()
	return nil
}

func (l *loginTest) iWillBeLoggedInSecurelyToMyAccount() error {
	if err := l.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	return nil
}

func (l *loginTest) theLoginShouldFailWith(msg string) error {
	if err := l.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := l.apiTest.AssertStringValueOnPathInResponse("error.message", msg); err != nil {
		return err
	}
	return nil
}

func (l *loginTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		l.apiTest.URL = "/v1/login"
		l.apiTest.Method = http.MethodPost
		l.apiTest.SetHeader("Content-Type", "application/json")
		l.apiTest.InitializeServer(l.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, err = l.DB.DeleteUser(ctx, l.user.ID)
		return ctx, err
	})

	ctx.Step(`^I am a registered user with details$`, l.iAmARegisteredUserWithDetails)
	ctx.Step(`^I fill the following details$`, l.iFillTheFollowingDetails)
	ctx.Step(`^I submit the registration form$`, l.iSubmitTheRegistrationForm)
	ctx.Step(`^I will be logged in securely to my account$`, l.iWillBeLoggedInSecurelyToMyAccount)
	ctx.Step(`^the login should fail with "([^"]*)"$`, l.theLoginShouldFailWith)
}
