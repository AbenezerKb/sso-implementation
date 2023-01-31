package change_password

import (
	"context"
	"net/http"
	"testing"

	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"

	"golang.org/x/crypto/bcrypt"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type forgetPasswordTest struct {
	test.TestInstance
	apiTest src.ApiTest
	User    db.User
	phone   string
}

func TestChangePassword(t *testing.T) {
	c := forgetPasswordTest{}
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeTest(t, "forget password", "features/forget_password.feature", c.InitializeScenario)
}

func (c *forgetPasswordTest) iHaveAUserAccountWithTheFollowingDetails(userDetails *godog.Table) error {
	userData, err := c.apiTest.ReadRow(userDetails, nil, false)
	if err != nil {
		return err
	}

	userValue := dto.User{}
	err = c.apiTest.UnmarshalJSON([]byte(userData), &userValue)
	if err != nil {
		return err
	}

	c.User, err = c.AuthenticateWithParam(userValue)
	if err != nil {
		return err
	}
	return nil
}

func (c *forgetPasswordTest) iFillMyPhoneNumberAs(phone string) error {
	c.phone = phone
	c.apiTest.URL = "/v1/resetCode"
	c.apiTest.Method = http.MethodGet
	c.apiTest.SetQueryParam("phone", phone)
	return nil
}

func (c *forgetPasswordTest) iRequestToHaveForgottenMyPassword() error {
	c.apiTest.SendRequest()
	return nil
}

func (c *forgetPasswordTest) iShouldSuccessfullyGetAChangePasswordRequestCode() error {
	return c.apiTest.AssertStatusCode(http.StatusOK)
}
func (c *forgetPasswordTest) iShouldSuccessfullyChangeMyPasswordUsingTheRequestCode() error {
	c.apiTest.ResetResponse()
	c.apiTest.Method = http.MethodPost
	c.apiTest.URL = "/v1/resetPassword"
	c.apiTest.SetBodyMap(map[string]interface{}{
		"phone":      c.phone,
		"password":   "somePass",
		"reset_code": "123455",
	})
	c.apiTest.SendRequest()

	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	fetchedUser, err := c.DB.GetUserById(context.Background(), c.User.ID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte("somePass")); err != nil {
		return err
	}

	return nil
}

func (c *forgetPasswordTest) iShouldFailChangeMyPasswordUsingAnIncorrectRequestCode() error {
	c.apiTest.ResetResponse()
	c.apiTest.Method = http.MethodPost
	c.apiTest.URL = "/v1/resetPassword"
	c.apiTest.SetBodyMap(map[string]interface{}{
		"phone":      c.phone,
		"password":   "somePass",
		"reset_code": "invalid",
	})
	c.apiTest.SendRequest()

	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.message", "invalid reset code")
}

func (c *forgetPasswordTest) InitializeScenario(ctx *godog.ScenarioContext) {

	c.apiTest.InitializeServer(c.Server)
	c.apiTest.SetHeader("Content-Type", "application/json")

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.DeleteUser(ctx, c.User.ID)
		return ctx, nil
	})

	ctx.Step(`^I fill my phone number as "([^"]*)"$`, c.iFillMyPhoneNumberAs)
	ctx.Step(`^I have a user account with the following details$`, c.iHaveAUserAccountWithTheFollowingDetails)
	ctx.Step(`^I request to have forgotten my password$`, c.iRequestToHaveForgottenMyPassword)
	ctx.Step(`^I should fail change my password using an incorrect request code$`, c.iShouldFailChangeMyPasswordUsingAnIncorrectRequestCode)
	ctx.Step(`^I should successfully change my password using the request code$`, c.iShouldSuccessfullyChangeMyPasswordUsingTheRequestCode)
	ctx.Step(`^I should successfully get a change password request code$`, c.iShouldSuccessfullyGetAChangePasswordRequestCode)
}
