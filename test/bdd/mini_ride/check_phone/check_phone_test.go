package check_phone

import (
	"context"
	"database/sql"
	"encoding/base64"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/middleware"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type checkPhoneTest struct {
	test.TestInstance
	apiTest    src.ApiTest
	Users      []dto.User
	checkPhone string
}

type checkPhoneRsp struct {
	Exists     string   `json:"exists"`
	CheckPhone string   `json:"check_phone"`
	User       dto.User `json:"user"`
}

func TestCheckPhone(t *testing.T) {
	c := &checkPhoneTest{}
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest = src.ApiTest{}
	c.apiTest.InitializeTest(t, "check phone for miniRide", "features/check_phone.feature", c.InitializeScenario)

}

func (c *checkPhoneTest) iAmAuthenticatedWithFollowingCredential(credential *godog.Table) error {
	credentialJson, err := c.apiTest.ReadRow(credential, nil, false)
	if err != nil {
		return err
	}
	var miniRideCredential middleware.MiniRideCredential
	err = c.apiTest.UnmarshalJSON([]byte(credentialJson), &miniRideCredential)
	if err != nil {
		return err
	}

	c.apiTest.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(miniRideCredential.UserName+":"+miniRideCredential.Password)))
	c.apiTest.SetHeader("Content-Type", "application/json")
	return nil
}
func (c *checkPhoneTest) theyAreTheFollowingUsersOnSso(users *godog.Table) error {
	usersJson, err := c.apiTest.ReadRows(users, nil, false)
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalJSON([]byte(usersJson), &c.Users)
	if err != nil {
		return err
	}

	for i := 0; i < len(c.Users); i++ {
		user, err := c.DB.CreateUserWithID(context.Background(), db.CreateUserWithIDParams{
			FirstName:      c.Users[i].FirstName,
			MiddleName:     c.Users[i].MiddleName,
			LastName:       c.Users[i].LastName,
			Phone:          c.Users[i].Phone,
			ProfilePicture: sql.NullString{String: c.Users[i].ProfilePicture, Valid: true},
			ID:             c.Users[i].ID,
		})
		if err != nil {
			return err
		}

		_, err = c.DB.UpdateUser(context.Background(), db.UpdateUserParams{
			Status: sql.NullString{String: c.Users[i].Status, Valid: true},
			ID:     user.ID,
		})
		if err != nil {
			return err
		}

		c.Users[i] = dto.User{
			ID:             user.ID,
			FirstName:      user.FirstName,
			MiddleName:     user.MiddleName,
			LastName:       user.LastName,
			Phone:          user.Phone,
			Status:         user.Status.String,
			ProfilePicture: user.ProfilePicture.String,
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checkPhoneTest) iRequestToCheckUsersWithTheFollowingPhone(checkPhone string) error {
	c.checkPhone = checkPhone
	c.apiTest.URL += checkPhone
	c.apiTest.SendRequest()
	return nil
}

func (c *checkPhoneTest) iShouldGetErrorWithMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	return c.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (c *checkPhoneTest) iShouldGetTheFollowingResponse(rsp *godog.Table) error {
	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	rspString, err := c.apiTest.ReadRow(rsp, []src.Type{
		{
			WithName: "user",
			Kind:     src.Object,
			Columns:  []string{"first_name", "middle_name", "last_name", "phone", "profile_picture", "status", "id"},
		},
	}, true)
	var checkRsp checkPhoneRsp
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalJSON([]byte(rspString), &checkRsp)
	if err != nil {
		return err
	}

	if checkRsp.Exists == "true" {
		fetchedUser, err := c.DB.GetUserByPhone(context.Background(), c.checkPhone)
		if err != nil {
			return err
		}

		if err := c.apiTest.AssertEqual(fetchedUser.FirstName, checkRsp.User.FirstName); err != nil {
			return err
		}
		if err := c.apiTest.AssertEqual(fetchedUser.MiddleName, checkRsp.User.MiddleName); err != nil {
			return err
		}
		if err := c.apiTest.AssertEqual(fetchedUser.LastName, checkRsp.User.LastName); err != nil {
			return err
		}
		if err := c.apiTest.AssertEqual(fetchedUser.Status.String, checkRsp.User.Status); err != nil {
			return err
		}
		if err := c.apiTest.AssertEqual(fetchedUser.Phone, checkRsp.User.Phone); err != nil {
			return err
		}
	}

	return nil
}

func (c *checkPhoneTest) InitializeScenario(ctx *godog.ScenarioContext) {

	c.apiTest.URL = "/v1/users/exists/"
	c.apiTest.Method = http.MethodGet
	c.apiTest.InitializeServer(c.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for i := 0; i < len(c.Users); i++ {
			_, _ = c.DB.DeleteUser(ctx, c.Users[i].ID)
		}
		return ctx, nil
	})

	ctx.Step(`^I request to check users with the following phone "([^"]*)"$`, c.iRequestToCheckUsersWithTheFollowingPhone)
	ctx.Step(`^I should get error with message "([^"]*)"$`, c.iShouldGetErrorWithMessage)
	ctx.Step(`^I should get the following response$`, c.iShouldGetTheFollowingResponse)
	ctx.Step(`^they are the following user\'s on sso$`, c.theyAreTheFollowingUsersOnSso)
	ctx.Step(`^I am authenticated with following credential$`, c.iAmAuthenticatedWithFollowingCredential)
}
