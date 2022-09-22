package get

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

type getProfileTest struct {
	test.TestInstance
	apiTest src.ApiTest
	user    db.User
}

func TestGetProfile(t *testing.T) {
	g := getProfileTest{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeTest(t, "get profile", "features/get_profile.feature", g.InitializeScenario)
}

func (g *getProfileTest) iAmLoggedInUserWithTheFollowingDetails(userCredentials *godog.Table) error {
	body, err := g.apiTest.ReadRow(userCredentials, nil, false)
	if err != nil {
		return err
	}

	userValue := dto.User{}
	err = g.apiTest.UnmarshalJSON([]byte(body), &userValue)
	if err != nil {
		return err
	}

	g.user, err = g.AuthenticateWithParam(userValue)
	if err != nil {
		return err
	}
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	return nil
}

func (g *getProfileTest) iRequestToGetMyProfile() error {
	g.apiTest.SendRequest()
	return nil
}

func (g *getProfileTest) iShouldSuccessfullyGetMyProfile() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	updatedUser := dto.User{}
	err := g.apiTest.UnmarshalResponseBodyPath("data", &updatedUser)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(updatedUser.Email, g.user.Email.String); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(updatedUser.FirstName, g.user.FirstName); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(updatedUser.LastName, g.user.LastName); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(updatedUser.Phone, g.user.Phone); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(updatedUser.ID, g.user.ID); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(updatedUser.CreatedAt, g.user.CreatedAt); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(updatedUser.Gender, g.user.Gender); err != nil {
		return err
	}

	return nil
}

func (g *getProfileTest) InitializeScenario(ctx *godog.ScenarioContext) {
	g.apiTest.URL = "/v1/profile"
	g.apiTest.Method = http.MethodGet
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.apiTest.InitializeServer(g.Server)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.user.ID)

		return ctx, nil
	})
	ctx.Step(`^I am logged in user with the following details$`, g.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I request to get my profile$`, g.iRequestToGetMyProfile)
	ctx.Step(`^I should successfully get my profile$`, g.iShouldSuccessfullyGetMyProfile)
}
