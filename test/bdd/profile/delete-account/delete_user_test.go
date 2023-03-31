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
	apiTest src.ApiTest
	user    db.User
}

func TestDeleteuser(t *testing.T) {

	c := &deleteuserTest{}
	c.TestInstance = test.Initiate("../../../../")

	c.apiTest.InitializeTest(t, "Delete user test", "features/delete_user.feature", c.InitializeScenario)
}
func (c *deleteuserTest) iAmLoggedInWithTheFollowingCreadentials(adminCredentials *godog.Table) error {
	var err error
	c.user, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, c.GrantRoleAfterFunc, err = c.GrantRoleForUserWithAfter(c.user.ID.String(), adminCredentials)
	if err != nil {
		return err
	}
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	return nil
}

func (c *deleteuserTest) iHaveAregistredAccountOnTheSystem(userForm *godog.Table) error {
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

func (c *deleteuserTest) iWantToDeleteMyAccount() error {
	c.apiTest.SendRequest()
	return nil
}

func (c *deleteuserTest) mYAccountShouldBeDeleted() error {
	if err := c.apiTest.AssertStatusCode(http.StatusNoContent); err != nil {
		return err
	}
	_, err := c.DB.GetUserById(context.Background(), c.user.ID)
	if err == nil {
		return fmt.Errorf("expected to not find the deleted users")
	}
	return nil
}
func (c *deleteuserTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		c.apiTest.URL = "/v1/profile"
		c.apiTest.Method = http.MethodDelete
		c.apiTest.SetHeader("Content-Type", "application/json")
		c.apiTest.InitializeServer(c.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, _ error) (context.Context, error) {
		//  delete the registered user
		_, _ = c.DB.DeleteUser(ctx, c.user.ID)
		// delete the admin
		return ctx, nil
	})

	ctx.Step(`^I have a registered account on the system$`, c.iHaveAregistredAccountOnTheSystem)
	ctx.Step(`^I am logged in with the following credentials$`, c.iAmLoggedInWithTheFollowingCreadentials)
	ctx.Step(`^I want to delete my account$`, c.iWantToDeleteMyAccount)
	ctx.Step(`^My account should be deleted$`, c.mYAccountShouldBeDeleted)

}
