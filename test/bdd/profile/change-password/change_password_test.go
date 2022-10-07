package change_password

import (
	"context"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type changePasswordTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	User        db.User
	newPassword string
}

func TestChangePassword(t *testing.T) {
	c := changePasswordTest{}
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeTest(t, "change password", "features/change_password.feature", c.InitializeScenario)
}

func (c *changePasswordTest) iAmLoggedInUserWithTheFollowingDetails(userDetails *godog.Table) error {
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
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	return nil
}

func (c *changePasswordTest) iFillTheFollowingDetails(changeInfo *godog.Table) error {
	newPassword, err := c.apiTest.ReadCell(changeInfo, "new_password", nil)
	if err != nil {
		return err
	}
	c.newPassword = fmt.Sprintf("%s", newPassword)

	body, err := c.apiTest.ReadRow(changeInfo, nil, false)
	if err != nil {
		return err
	}

	c.apiTest.Body = body
	return nil
}

func (c *changePasswordTest) iRequestToChangeMyPassword() error {
	c.apiTest.SendRequest()
	return nil
}

func (c *changePasswordTest) iShouldSuccessfullyChangeMyPassword() error {
	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	fetchedUser, err := c.DB.GetUserById(context.Background(), c.User.ID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(c.newPassword)); err != nil {
		return err
	}

	return nil
}

func (c *changePasswordTest) thePasswordChangingShouldFailWithFieldErrorMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (c *changePasswordTest) thePasswordChangingShouldFailWithMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (c *changePasswordTest) InitializeScenario(ctx *godog.ScenarioContext) {

	c.apiTest.URL = "/v1/profile/password"
	c.apiTest.Method = http.MethodPatch
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.SetHeader("Content-Type", "application/json")

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.DeleteUser(ctx, c.User.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in user with the following details$`, c.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I fill the following details$`, c.iFillTheFollowingDetails)
	ctx.Step(`^I request to change my password$`, c.iRequestToChangeMyPassword)
	ctx.Step(`^I should successfully change my password$`, c.iShouldSuccessfullyChangeMyPassword)
	ctx.Step(`^The password changing should fail with field error message "([^"]*)"$`, c.thePasswordChangingShouldFailWithFieldErrorMessage)
	ctx.Step(`^The password changing should fail with message "([^"]*)"$`, c.thePasswordChangingShouldFailWithMessage)
}
