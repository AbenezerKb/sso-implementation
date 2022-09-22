package update_user_status

import (
	"context"
	"database/sql"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type updateUserStatusTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Admin   db.User
	user    db.User
}

func TestUpdateUserStatus(t *testing.T) {
	u := updateUserStatusTest{}
	u.TestInstance = test.Initiate("../../../../")

	u.apiTest.URL = "/v1/users/"
	u.apiTest.Method = http.MethodPatch
	u.apiTest.InitializeServer(u.Server)

	u.apiTest.InitializeTest(t, "update user status", "features/update_user_status.feature", u.InitializeScenario)
}

func (u *updateUserStatusTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
	body, err := u.apiTest.ReadRow(adminCredential, nil, false)
	if err != nil {
		return err
	}

	adminValue := dto.User{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &adminValue)
	if err != nil {
		return err
	}

	u.Admin, err = u.AuthenticateWithParam(adminValue)
	if err != nil {
		return err
	}
	u.apiTest.SetHeader("Authorization", "Bearer "+u.AccessToken)
	return u.GrantRoleForUser(u.Admin.ID.String(), adminCredential)
}

func (u *updateUserStatusTest) iUpdateTheUsersStatusTo(status string) error {
	u.apiTest.SetBodyValue("status", status)
	u.apiTest.SendRequest()
	return nil
}

func (u *updateUserStatusTest) theUserStatusShouldUpdateTo(status string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	updatedUserData, err := u.DB.GetUserById(context.Background(), u.user.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(updatedUserData.Status, status); err != nil {
		return err
	}
	return nil
}

func (u *updateUserStatusTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return u.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message)
}

func (u *updateUserStatusTest) thereIsUserWithId(userID string) error {
	u.apiTest.URL += userID + "/status"
	return nil
}

func (u *updateUserStatusTest) thereIsUserWithTheFollowingDetails(userDetails *godog.Table) error {
	body, err := u.apiTest.ReadRow(userDetails, nil, false)

	if err != nil {
		return err
	}

	userValues := dto.User{}
	err = u.apiTest.UnmarshalJSON([]byte(body), &userValues)
	if err != nil {
		return err
	}
	u.user, err = u.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:      userValues.FirstName,
		MiddleName:     userValues.MiddleName,
		LastName:       userValues.LastName,
		Email:          sql.NullString{String: userValues.Email, Valid: true},
		Phone:          userValues.Phone,
		UserName:       userValues.UserName,
		Gender:         userValues.Gender,
		ProfilePicture: sql.NullString{String: userValues.ProfilePicture, Valid: true},
	})
	u.apiTest.URL += u.user.ID.String() + "/status"
	return err
}

func (u *updateUserStatusTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Step(`^I am logged in as admin user$`, u.iAmLoggedInAsAdminUser)
	ctx.Step(`^I update the user\'s status to "([^"]*)"$`, u.iUpdateTheUsersStatusTo)
	ctx.Step(`^the user status should update to "([^"]*)"$`, u.theUserStatusShouldUpdateTo)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, u.thenIShouldGetErrorWithMessage)
	ctx.Step(`^there is user with id "([^"]*)"$`, u.thereIsUserWithId)
	ctx.Step(`^there is user with the following details:$`, u.thereIsUserWithTheFollowingDetails)
}
