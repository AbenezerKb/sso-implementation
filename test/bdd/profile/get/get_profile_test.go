package get

import (
	"context"
	"fmt"
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
	apiTest  src.ApiTest
	user     db.User
	userRole string
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
	g.Conn.Query(context.Background(), fmt.Sprintf("insert into casbin_rule (p_type, v0, v1) values('g','%s', '%s');", g.user.ID, userValue.Role))

	g.userRole = userValue.Role

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
	fetchedUser := dto.User{}
	err := g.apiTest.UnmarshalResponseBodyPath("data", &fetchedUser)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedUser.Email, g.user.Email.String); err != nil {
		return err
	}
	if err := g.apiTest.AssertEqual(fetchedUser.FirstName, g.user.FirstName); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedUser.LastName, g.user.LastName); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedUser.Phone, g.user.Phone); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedUser.ID, g.user.ID); err != nil {
		return err
	}

	if equal := fetchedUser.CreatedAt.Equal(g.user.CreatedAt); !equal {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedUser.Gender, g.user.Gender); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(fetchedUser.Role, g.userRole); err != nil {
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
