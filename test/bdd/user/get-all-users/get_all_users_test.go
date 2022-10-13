package get_all_users

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
)

type getUsersTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	users       []db.User
	Admin       db.User
	Preferences preferenceData
}

type preferenceData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func TestGetUsers(t *testing.T) {
	c := getUsersTest{}
	c.apiTest.URL = "/v1/users"
	c.apiTest.Method = http.MethodGet
	c.TestInstance = test.Initiate("../../../../")
	c.apiTest.InitializeServer(c.Server)
	c.apiTest.InitializeTest(t, "get users test", "features/get_all_users.feature", c.InitializeScenario)
}

func (c *getUsersTest) theFollowingUsersAreRegisteredOnTheSystem(users *godog.Table) error {
	usersJSON, err := c.apiTest.ReadRows(users, nil, false)
	if err != nil {
		return err
	}
	var usersData []dto.User
	err = c.apiTest.UnmarshalJSON([]byte(usersJSON), &usersData)
	if err != nil {
		return err
	}
	for _, v := range usersData {
		user, err := c.DB.CreateUser(context.Background(), db.CreateUserParams{
			FirstName:      v.FirstName,
			MiddleName:     v.MiddleName,
			LastName:       v.LastName,
			Email:          sql.NullString{String: v.Email, Valid: true},
			Phone:          v.Phone,
			UserName:       v.UserName,
			Gender:         v.Gender,
			ProfilePicture: sql.NullString{String: v.ProfilePicture, Valid: true},
		})
		if err != nil {
			return err
		}
		c.users = append(c.users, user)
	}

	return nil
}

func (c *getUsersTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	c.Admin, err = c.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, c.GrantRoleAfterFunc, err = c.GrantRoleForUserWithAfter(c.Admin.ID.String(), adminCredentials)
	return err
}

func (c *getUsersTest) iRequestToGetAllTheUsersWithTheFollowingPreferences(preferences *godog.Table) error {
	preferencesJSON, err := c.apiTest.ReadRow(preferences, []src.Type{
		{
			Column: "page",
			Kind:   src.Any,
		},
		{
			Column: "per_page",
			Kind:   src.Any,
		},
	}, false)
	if err != nil {
		return err
	}
	err = c.apiTest.UnmarshalJSON([]byte(preferencesJSON), &c.Preferences)
	if err != nil {
		return err
	}

	c.apiTest.SetQueryParam("page", fmt.Sprintf("%d", c.Preferences.Page))
	c.apiTest.SetQueryParam("per_page", fmt.Sprintf("%d", c.Preferences.PerPage))
	c.apiTest.SetHeader("Authorization", "Bearer "+c.AccessToken)
	c.apiTest.SendRequest()
	return nil
}

func (c *getUsersTest) iShouldGetTheListOfUsersThatPassMyPreferences() error {
	var responseUsers []dto.User
	var metaData model.MetaData

	if err := c.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	err := c.apiTest.UnmarshalResponseBodyPath("meta_data", &metaData)
	if err != nil {
		return err
	}

	err = c.apiTest.UnmarshalResponseBodyPath("data", &responseUsers)
	if err != nil {
		return err
	}
	var total int
	if c.Preferences.Page < metaData.Total/c.Preferences.PerPage {
		total = c.Preferences.PerPage
	} else if c.Preferences.Page == metaData.Total/c.Preferences.PerPage {
		total = metaData.Total % c.Preferences.PerPage
	} else {
		total = 0
	}
	if err := c.apiTest.AssertEqual(len(responseUsers), total); err != nil {
		return err
	}
	for _, v := range responseUsers {
		found := false
		for _, v2 := range append(c.users, c.Admin) {
			if v.ID.String() == v2.ID.String() {
				found = true
				continue
			}
		}
		if !found {
			return fmt.Errorf("expected user: %v", v)
		}
	}
	return nil
}

func (c *getUsersTest) iShouldGetErrorMessage(message string) error {
	if err := c.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := c.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (c *getUsersTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		for _, v := range c.users {
			_, _ = c.DB.DeleteUser(ctx, v.ID)
		}
		_, _ = c.DB.DeleteUser(ctx, c.Admin.ID)
		return ctx, nil
	})
	ctx.Step(`^I am logged in as admin user$`, c.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get all the users with the following preferences$`, c.iRequestToGetAllTheUsersWithTheFollowingPreferences)
	ctx.Step(`^I should get the list of users that pass my preferences$`, c.iShouldGetTheListOfUsersThatPassMyPreferences)
	ctx.Step(`^I should get error message "([^"]*)"$`, c.iShouldGetErrorMessage)
	ctx.Step(`^The following users are registered on the system$`, c.theFollowingUsersAreRegisteredOnTheSystem)
}
