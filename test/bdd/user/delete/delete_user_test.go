package deleteUser

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
)

type deleteuserTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	admin, user db.User
}

func TestDeleteuser(t *testing.T) {

	c := &deleteuserTest{}
	c.TestInstance = test.Initiate("../../../../")

	c.apiTest.InitializeTest(t, "Delete user test", "features/delete_user.feature", c.InitializeScenario)
}
func (c *deleteuserTest) iAmLoggedInWithTheFollowingCreadentials(adminCredentials *godog.Table) error {
	var err error
	c.admin, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, c.GrantRoleAfterFunc, err = c.GrantRoleForUserWithAfter(c.admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	return nil
}

func (c *deleteuserTest) iHaveAregistredUsers(userForm *godog.Table) error {
	body, err := c.apiTest.ReadRow(userForm, nil, false)
	if err != nil {
		return err
	}
	var user dto.User
	err = c.apiTest.UnmarshalJSON([]byte(body), &user)
	if err != nil {
		return err
	}
	c.user, err = c.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email: sql.NullString{
			Valid:  true,
			String: user.Email,
		},
		Phone:          user.Phone,
		UserName:       user.UserName,
		Password:       user.Password,
		Gender:         user.Gender,
		ProfilePicture: sql.NullString{String: user.ProfilePicture, Valid: true}})
	if err != nil {
		return err
	}

	return nil
}

func (c *deleteuserTest) iRequestToDeleteTheUser() error {
	c.apiTest.URL = c.apiTest.URL + c.user.ID.String()
	c.apiTest.SendRequest()
	return nil
}

func (c *deleteuserTest) theUserShouldBeDeleted() error {
	if err := c.apiTest.AssertStatusCode(http.StatusNoContent); err != nil {
		return err
	}
	_, err := c.DB.GetUserById(context.Background(), c.user.ID)
	if err == nil {
		return fmt.Errorf("expected to not find the deleted users")
	}
	return nil
}
func (c *deleteuserTest) iRequestToDeleteTheUsersWithIn(userID string) error {
	c.apiTest.URL = c.apiTest.URL + userID
	c.apiTest.SendRequest()
	return nil
}
func (c *deleteuserTest) theSystemUserShouldGetAnErrorMessage(message string) error {
	err := c.apiTest.AssertStatusCode(http.StatusBadRequest)
	if err != nil {
		err := c.apiTest.AssertStatusCode(http.StatusNotFound)
		if err != nil {
			err := c.apiTest.AssertStatusCode(http.StatusInternalServerError)
			return err
		}
	}
	err = c.apiTest.AssertStringValueOnPathInResponse("error.message", message)
	if err != nil {
		return err
	}
	return nil
}

func (c *deleteuserTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/users/"
		c.apiTest.Method = http.MethodDelete
		c.apiTest.SetHeader("Content-Type", "application/json")
		c.apiTest.InitializeServer(c.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, _ error) (context.Context, error) {
		//  delete the registered user
		_, _ = c.DB.DeleteUser(ctx, c.user.ID)
		// delete the admin
		_, _ = c.DB.DeleteUser(ctx, c.admin.ID)
		return ctx, nil
	})

	ctx.Step(`^I am logged in with the following credentials$`, c.iAmLoggedInWithTheFollowingCreadentials)
	ctx.Step(`^I have a registered users$`, c.iHaveAregistredUsers)
	ctx.Step(`^I request to delete the user$`, c.iRequestToDeleteTheUser)
	ctx.Step(`^the user should be deleted$`, c.theUserShouldBeDeleted)
	ctx.Step(`^I request to delete the users with in "([^"]*)"$`,
		c.iRequestToDeleteTheUsersWithIn)
	ctx.Step(`^The system user should get an error message "([^"]*)"$`,
		c.theSystemUserShouldGetAnErrorMessage)

}
