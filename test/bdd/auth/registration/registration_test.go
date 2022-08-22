package registration

import (
	"context"
	"encoding/json"
	"github.com/dongri/phonenumber"
	"net/http"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type registrationTest struct {
	test.TestInstance
	apiTest src.ApiTest
	user    struct {
		OK   bool     `json:"ok"`
		Data dto.User `json:"data"`
	}
	redisSeeder seed.RedisDB
}

func TestRegistertion(t *testing.T) {

	a := &registrationTest{}
	a.TestInstance = test.Initiate("../../../../")
	a.redisSeeder = seed.RedisDB{DB: a.Redis}

	a.apiTest.InitializeTest(t, "Login test", "features/registration.feature", a.InitializeScenario)
}

func (r *registrationTest) iFillTheFormWithTheFollowingDetails(userForm *godog.Table) error {

	// set otp to redis
	phone, err := r.apiTest.ReadCellString(userForm, "phone")
	if err != nil {
		return err
	}
	otp, err := r.apiTest.ReadCellString(userForm, "otp")
	if err != nil {
		return err
	}
	err = r.redisSeeder.Feed(seed.RedisModel{
		Key:   phonenumber.Parse(phone, "ET"),
		Value: otp,
	})
	if err != nil {
		return err
	}

	body, err := r.apiTest.ReadRow(userForm, nil, false)
	if err != nil {
		return err
	}
	r.apiTest.Body = body
	return nil
}

func (r *registrationTest) iSubmitTheRegistrationForm() error {
	r.apiTest.SendRequest()
	return nil
}

func (r *registrationTest) iWillHaveANewAccount() error {
	if err := r.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	err := json.Unmarshal(r.apiTest.ResponseBody, &r.user)
	return err
}

func (r *registrationTest) theRegistrationShouldFailWith(msg string) error {
	if err := r.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if err := r.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", msg); err != nil {
		return err
	}

	return nil
}

func (r *registrationTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {

		r.apiTest.URL = "/v1/register"
		r.apiTest.Method = http.MethodPost
		r.apiTest.SetHeader("Content-Type", "application/json")
		r.apiTest.InitializeServer(r.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(ctx, r.user.Data.ID)
		return ctx, nil
	})

	ctx.Step(`^I fill the form with the following details$`, r.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I submit the registration form$`, r.iSubmitTheRegistrationForm)
	ctx.Step(`^I will have a new account$`, r.iWillHaveANewAccount)
	ctx.Step(`^the registration should fail with "([^"]*)"$`, r.theRegistrationShouldFailWith)
}
