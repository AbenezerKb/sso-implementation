package get_current_sessions

import (
	"context"
	"fmt"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getCurrentSessionsTest struct {
	test.TestInstance
	apiTest         src.ApiTest
	user            db.User
	currentSessions []db.Internalrefreshtoken
}

func TestGetCurrentSessions(t *testing.T) {
	g := getCurrentSessionsTest{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeTest(t, "get current sessions", "features/get_current_sessions.feature", g.InitializeScenario)
}

func (g *getCurrentSessionsTest) iAmLoggedInUserWithTheFollowingDetails(userCredentials *godog.Table) error {
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

func (g *getCurrentSessionsTest) andIHaveTheFollowingSessionsOnTheSystem(sessions *godog.Table) error {
	sessionsJSON, err := g.apiTest.ReadRows(sessions, nil, false)
	if err != nil {
		return err
	}

	var sessionsData []dto.InternalRefreshToken
	err = g.apiTest.UnmarshalJSON([]byte(sessionsJSON), &sessionsData)
	if err != nil {
		return err
	}
	for _, v := range sessionsData {
		session, err := g.DB.SaveInternalRefreshToken(context.Background(), db.SaveInternalRefreshTokenParams{
			ExpiresAt:    time.Now().Add(time.Hour * 2),
			UserID:       g.user.ID,
			RefreshToken: v.RefreshToken,
			IpAddress:    v.IPAddress,
			UserAgent:    v.UserAgent,
		})
		if err != nil {
			return err
		}
		g.currentSessions = append(g.currentSessions, session)
	}

	return nil
}

func (g *getCurrentSessionsTest) iRequestToGetMyCurrentSessions() error {
	g.apiTest.SendRequest()
	return nil
}

func (g *getCurrentSessionsTest) iShouldGetTheAllMySessions() error {
	var responseSessions []dto.InternalRefreshToken

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	err := g.apiTest.UnmarshalResponseBodyPath("data", &responseSessions)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(len(responseSessions), len(g.currentSessions)+1); err != nil {
		return err
	}

	for _, v := range g.currentSessions {
		found := false
		for _, v2 := range responseSessions {
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

func (g *getCurrentSessionsTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		g.apiTest.URL = "/v1/profile/devices"
		g.apiTest.Method = http.MethodGet
		g.apiTest.SetHeader("Content-Type", "application/json")
		g.apiTest.InitializeServer(g.Server)

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.user.ID)

		_ = g.DB.RemoveInternalRefreshTokenByUserID(ctx, g.user.ID)

		return ctx, nil
	})
	ctx.Step(`^And I have the following sessions on the system$`, g.andIHaveTheFollowingSessionsOnTheSystem)
	ctx.Step(`^I am logged in user with the following details$`, g.iAmLoggedInUserWithTheFollowingDetails)
	ctx.Step(`^I request to get my current sessions$`, g.iRequestToGetMyCurrentSessions)
	ctx.Step(`^I should get the all my sessions$`, g.iShouldGetTheAllMySessions)
}
