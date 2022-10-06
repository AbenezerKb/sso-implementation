package change_phone

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src/seed"
)

type changePhoneTest struct {
	test.TestInstance
	apiTest           src.ApiTest
	AuthenticatedUser db.User
	User              db.User
	redisSeeder       seed.RedisDB
	newPhone          string
}

func TestChangePhone(t *testing.T) {
	c := changePhoneTest{}
	c.TestInstance = test.Initiate("../../../../")
	// create redis seeder
	c.redisSeeder = seed.RedisDB{DB: c.Redis}
	c.apiTest.InitializeTest(t, "change phone", "features/change_phone.feature", c.InitializeScenario)
}

func (c *changePhoneTest) iAmLoggedInUserWithTheFollowingDetails(userDetails *godog.Table) error {
	userData, err := c.apiTest.ReadRow(userDetails, nil, false)
	if err != nil {
		return err
	}

	userValue := dto.User{}
	err = c.apiTest.UnmarshalJSON([]byte(userData), &userValue)
	if err != nil {
		return err
	}

	c.AuthenticatedUser, err = c.AuthenticateWithParam(userValue)
	if err != nil {
		return err
	}
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	return nil
}

func (c *changePhoneTest) theFollowingUserIsRegisteredOnTheSystem(userDetails *godog.Table) error {
	body, err := c.apiTest.ReadRow(userDetails, nil, false)

	if err != nil {
		return err
	}

	userValues := dto.User{}
	err = c.apiTest.UnmarshalJSON([]byte(body), &userValues)
	if err != nil {
		return err
	}
	c.User, err = c.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:      userValues.FirstName,
		MiddleName:     userValues.MiddleName,
		LastName:       userValues.LastName,
		Email:          sql.NullString{String: userValues.Email, Valid: true},
		Phone:          userValues.Phone,
		UserName:       userValues.UserName,
		Gender:         userValues.Gender,
		ProfilePicture: sql.NullString{String: userValues.ProfilePicture, Valid: true},
	})

	return err
}

func (c *changePhoneTest) iFillTheFollowingDetails(changeInfo *godog.Table) error {
	// set otp to redis
	phone, err := c.apiTest.ReadCell(changeInfo, "phone", nil)
	if err != nil {
		return err
	}

	otp, err := c.apiTest.ReadCell(changeInfo, "otp", nil)
	if err != nil {
		return err
	}

	err = c.redisSeeder.Feed(seed.RedisModel{
		Key:   fmt.Sprintf("%s", phone),
		Value: fmt.Sprintf("%s", otp),
	})
	if err != nil {
		return err
	}

	c.newPhone = fmt.Sprintf("%s", phone)

	body, err := c.apiTest.ReadRow(changeInfo, nil, false)
	if err != nil {
		return err
	}

	c.apiTest.Body = body
	return nil
}

func (c *changePhoneTest) iFillTheFollowingDetailsWithWrongInfo(changeInfo *godog.Table) error {
	body, err := c.apiTest.ReadRow(changeInfo, nil, false)
	if err != nil {
		return err
	}
	c.apiTest.Body = body
	return nil
}

func (c *changePhoneTest) iRequestToChangeMyPhone() error {
	c.apiTest.SendRequest()
	return nil
}

func (c *changePhoneTest) iShouldSuccessfullyChangeMyPhone() error {
	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	fetchedUser, err := c.DB.GetUserById(context.Background(), c.AuthenticatedUser.ID)
	if err != nil {
		return err
	}

	if err := c.apiTest.AssertEqual(fetchedUser.Phone, c.newPhone); err != nil {
		return err
	}

	return nil
}

func (c *changePhoneTest) thePhoneChangingShouldFailWithMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (c *changePhoneTest) InitializeScenario(ctx *godog.ScenarioContext) {
	c.apiTest.URL = "/v1/profile/phone"
	c.apiTest.Method = http.MethodPatch
	c.apiTest.InitializeServer(c.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = c.DB.DeleteUser(ctx, c.User.ID)
		_, _ = c.DB.DeleteUser(ctx, c.AuthenticatedUser.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in user with the following details$`, c.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I fill the following details$`, c.iFillTheFollowingDetails)
	ctx.Step(`^I request to change my phone$`, c.iRequestToChangeMyPhone)
	ctx.Step(`^I should successfully change my phone$`, c.iShouldSuccessfullyChangeMyPhone)
	ctx.Step(`^The following user is registered on the system$`, c.theFollowingUserIsRegisteredOnTheSystem)
	ctx.Step(`^The phone changing should fail with message "([^"]*)"$`, c.thePhoneChangingShouldFailWithMessage)
	ctx.Step(`^I fill the following details with wrong info$`, c.iFillTheFollowingDetailsWithWrongInfo)
}
