package get

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

type getUserTest struct {
	test.TestInstance
	apiTest src.ApiTest
	Admin   db.User
	user    db.User
}

func TestGetUser(t *testing.T) {
	g := getUserTest{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeTest(t, "get user", "features/get_user.feature", g.InitializeScenario)
}

func (g *getUserTest) iAmLoggedInAsAdminUser(adminCredential *godog.Table) error {
	body, err := g.apiTest.ReadRow(adminCredential, nil, false)
	if err != nil {
		return err
	}

	adminValue := dto.User{}
	err = g.apiTest.UnmarshalJSON([]byte(body), &adminValue)
	if err != nil {
		return err
	}

	g.Admin, err = g.AuthenticateWithParam(adminValue)
	if err != nil {
		return err
	}
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	return g.GrantRoleForUser(g.Admin.ID.String(), adminCredential)
}

func (g *getUserTest) thereIsUserWithTheFollowingDetails(userDetails *godog.Table) error {
	body, err := g.apiTest.ReadRow(userDetails, nil, false)

	if err != nil {
		return err
	}

	userValues := dto.User{}
	err = g.apiTest.UnmarshalJSON([]byte(body), &userValues)
	if err != nil {
		return err
	}
	g.user, err = g.DB.CreateUser(context.Background(), db.CreateUserParams{
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

func (g *getUserTest) iHaveUsersId() error {
	g.apiTest.URL += g.user.ID.String()
	return nil
}

func (g *getUserTest) iHaveUserWithId(userID string) error {
	g.apiTest.URL += userID
	return nil
}

func (g *getUserTest) iGetTheUser() error {
	g.apiTest.SendRequest()
	return nil
}

func (g *getUserTest) iShouldSuccessfullyGetTheUser() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	getedUser := dto.User{}
	err := g.apiTest.UnmarshalResponseBodyPath("data", &getedUser)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(getedUser.Email, g.user.Email.String); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(getedUser.FirstName, g.user.FirstName); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(getedUser.LastName, g.user.LastName); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(getedUser.Phone, g.user.Phone); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(getedUser.ID, g.user.ID); err != nil {
		return err
	}

	return nil
}

func (g *getUserTest) thenIShouldGetErrorWithMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return g.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}
func (g *getUserTest) InitializeScenario(ctx *godog.ScenarioContext) {
	g.apiTest.URL = "/v1/users/"
	g.apiTest.Method = http.MethodGet
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.apiTest.InitializeServer(g.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.Admin.ID)
		_, _ = g.DB.DeleteUser(ctx, g.user.ID)

		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I Get the user$`, g.iGetTheUser)
	ctx.Step(`^I have user with id "([^"]*)"$`, g.iHaveUserWithId)
	ctx.Step(`^I have users id$`, g.iHaveUsersId)
	ctx.Step(`^I should successfully get the user$`, g.iShouldSuccessfullyGetTheUser)
	ctx.Step(`^Then I should get error with message "([^"]*)"$`, g.thenIShouldGetErrorWithMessage)
	ctx.Step(`^there is user with the following details:$`, g.thereIsUserWithTheFollowingDetails)
}
