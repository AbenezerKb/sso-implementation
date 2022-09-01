package update

import (
	"context"
	"database/sql"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"github.com/dongri/phonenumber"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateUserTest struct {
	test.TestInstance
	User          db.User
	NewUserData   dto.User
	apiTest       src.ApiTest
	DublicateUser dto.User
}

func TestUpdateUser(t *testing.T) {
	u := updateUserTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.InitializeTest(t, "update user", "features/update_user.feature", u.InitializeScenario)
}

func (u *updateUserTest) iAmLoggedInUserWithTheFollowingDetails(userDetails *godog.Table) error {
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

func (u *updateUserTest) iFillTheFormWithTheFollowingDetails(updateDataArg *godog.Table) error {
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

func (u *updateUserTest) iUpdateMyProfile() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *updateUserTest) myProfileShouldBeUpdated() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	updatedUserData, err := u.DB.GetUserById(context.Background(), u.User.ID)
	if err != nil {
		return err
	}

	if u.NewUserData.Email != "" {
		if err := u.apiTest.AssertEqual(u.NewUserData.Email, updatedUserData.Email.String); err != nil {
			return err
		}
	}

	if u.NewUserData.FirstName != "" {
		if err := u.apiTest.AssertEqual(u.NewUserData.FirstName, updatedUserData.FirstName); err != nil {
			return err
		}
	}

	if u.NewUserData.MiddleName != "" {
		if err := u.apiTest.AssertEqual(u.NewUserData.MiddleName, updatedUserData.MiddleName); err != nil {
			return err
		}
	}

	if u.NewUserData.LastName != "" {
		if err := u.apiTest.AssertEqual(u.NewUserData.LastName, updatedUserData.LastName); err != nil {
			return err
		}
	}

	if u.NewUserData.Phone != "" {
		newPhone := phonenumber.Parse(u.NewUserData.Phone, "ET")

		if err := u.apiTest.AssertEqual(newPhone, updatedUserData.Phone); err != nil {
			return err
		}
	}

	return nil
}
func (u *updateUserTest) thereIsUserWithFollowingDetails(dublicateUser *godog.Table) error {
	body, err := u.apiTest.ReadRow(dublicateUser, nil, false)
	if err != nil {
		return err
	}
	if err := u.apiTest.UnmarshalJSONAt([]byte(body), "", &u.DublicateUser); err != nil {
		return err
	}

	userData, err := u.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:  u.DublicateUser.FirstName,
		Phone:      u.DublicateUser.Phone,
		Email:      sql.NullString{String: u.DublicateUser.Email, Valid: true},
		MiddleName: u.DublicateUser.MiddleName,
		LastName:   u.DublicateUser.LastName,
	})
	if err != nil {
		return err
	}
	u.DublicateUser.ID = userData.ID
	return nil

}

func (u *updateUserTest) theUpdateShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (u *updateUserTest) theUpdateShouldFailWithErrorMessage(message string) error {

	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (u *updateUserTest) InitializeScenario(ctx *godog.ScenarioContext) {
	u.apiTest.URL = "/v1/users"
	u.apiTest.Method = http.MethodPatch
	u.apiTest.SetHeader("Content-Type", "application/json")
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteUser(ctx, u.User.ID)
		_, _ = u.DB.DeleteUser(ctx, u.DublicateUser.ID)

		return ctx, nil
	})

	ctx.Step(`^I am logged in user with the following details:$`, u.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I fill the form with the following details:$`, u.iFillTheFormWithTheFollowingDetails)
	ctx.Step(`^I update my profile$`, u.iUpdateMyProfile)
	ctx.Step(`^my profile should be updated$`, u.myProfileShouldBeUpdated)
	ctx.Step(`^The update should fail with message "([^"]*)"$`, u.theUpdateShouldFailWithMessage)
	ctx.Step(`^there is user with following details:$`, u.thereIsUserWithFollowingDetails)
	ctx.Step(`^The update should fail with error message "([^"]*)"$`, u.theUpdateShouldFailWithErrorMessage)

}
