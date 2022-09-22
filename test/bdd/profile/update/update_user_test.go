package update

import (
	"context"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateUserProfileTest struct {
	test.TestInstance
	User        db.User
	NewUserData dto.User
	apiTest     src.ApiTest
}

func TestUpdateProfile(t *testing.T) {
	u := updateUserProfileTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.InitializeTest(t, "update user", "features/update_user.feature", u.InitializeScenario)
}

func (u *updateUserProfileTest) iAmLoggedInUserWithTheFollowingDetails(userDetails *godog.Table) error {
	userData, err := u.apiTest.ReadRow(userDetails, nil, false)
	if err != nil {
		return err
	}

	userValue := dto.User{}
	err = u.apiTest.UnmarshalJSON([]byte(userData), &userValue)
	if err != nil {
		return err
	}

	u.User, err = u.AuthenticateWithParam(userValue)
	if err != nil {
		return err
	}
	u.apiTest.SetHeader("Authorization", "Bearer "+u.AccessToken)
	return nil
}

func (u *updateUserProfileTest) iFillTheFormWithTheFollowingDetails(updateDataArg *godog.Table) error {
	userData, err := u.apiTest.ReadRow(updateDataArg, nil, false)
	if err != nil {
		return err
	}

	err = u.apiTest.UnmarshalJSON([]byte(userData), &u.NewUserData)
	if err != nil {
		return err
	}

	u.apiTest.Body = userData
	return nil
}

func (u *updateUserProfileTest) iUpdateMyProfile() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *updateUserProfileTest) myProfileShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	updatedUserData, err := u.DB.GetUserById(context.Background(), u.User.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(u.NewUserData.FirstName, updatedUserData.FirstName); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(u.NewUserData.MiddleName, updatedUserData.MiddleName); err != nil {
		return err
	}
	if err := u.apiTest.AssertEqual(u.NewUserData.LastName, updatedUserData.LastName); err != nil {
		return err
	}

	return nil
}

func (u *updateUserProfileTest) theUpdateShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (u *updateUserProfileTest) InitializeScenario(ctx *godog.ScenarioContext) {
	u.apiTest.URL = "/v1/profile"
	u.apiTest.Method = http.MethodPut
	u.apiTest.SetHeader("Content-Type", "application/json")
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteUser(ctx, u.User.ID)

		return ctx, nil
	})

	ctx.Step(`^I am logged in user with the following details:$`, u.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I fill the form with the following details:$`, u.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I update my profile$`, u.iUpdateMyProfile)
	ctx.Step(`^my profile should be updated$`, u.myProfileShouldBeUpdated)
	ctx.Step(`^The update should fail with message "([^"]*)"$`, u.theUpdateShouldFailWithMessage)

}
