package reset_user_password

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type resetUserPasswordTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Admin   db.User
	user    db.User
}

func TestUpdateUserStatus(t *testing.T) {
	u := resetUserPasswordTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.InitializeTest(t, "reset user password", "features/reset_user_password.feature", u.InitializeScenario)
}

func (u *resetUserPasswordTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
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

func (u *resetUserPasswordTest) thereIsUserWithId(userID string) error {
	u.apiTest.URL += userID + "/status"
	return nil
}

func (u *resetUserPasswordTest) thereIsUserWithTheFollowingDetails(userDetails *godog.Table) error {
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
		Password:       "somePass",
		ProfilePicture: sql.NullString{String: userValues.ProfilePicture, Valid: true},
	})

	u.apiTest.URL += u.user.ID.String() + "/password"
	return err
}

func (u *resetUserPasswordTest) iResetTheUsersPassword() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *resetUserPasswordTest) theUsersPasswordShouldBeChanged() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	user, err := u.DB.GetUserById(context.Background(), u.user.ID)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(user.Password, u.user.Password); err == nil {
		return fmt.Errorf("password was not changed")
	}
	return nil
}

func (u *resetUserPasswordTest) InitializeScenario(ctx *godog.ScenarioContext) {

	u.apiTest.URL = "/v1/users/"
	u.apiTest.Method = http.MethodPatch
	u.apiTest.InitializeServer(u.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteUser(ctx, u.user.ID)
		_, _ = u.DB.DeleteUser(ctx, u.Admin.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, u.iAmLoggedInAsAdminUser)
	ctx.Step(`^there is user with the following details:$`, u.thereIsUserWithTheFollowingDetails)
	ctx.Step(`^I reset the user\'s password"$`, u.iResetTheUsersPassword)
	ctx.Step(`^the user\'s password should be changed$`, u.theUsersPasswordShouldBeChanged)

}
