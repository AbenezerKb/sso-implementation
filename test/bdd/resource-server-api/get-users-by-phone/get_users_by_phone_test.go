package get_users_by_phone

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils/collection"
	"sso/test"
	"testing"
)

type getUsersByIDOrPhone struct {
	test.TestInstance
	apiTest        src.ApiTest
	users          []db.User
	resourceServer db.ResourceServer
}

func TestGetUsersByIDOrPhone(t *testing.T) {
	g := getUsersByIDOrPhone{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.Server = g.Server
	g.apiTest.URL = "/v1/internal/users"
	g.apiTest.Method = http.MethodPost
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.apiTest.RunTest(t,
		"get users by id or phone test",
		&src.TestOptions{
			Paths: []string{
				"features/get_users_by_phone.feature",
				"features/get_users_by_id.feature",
			},
		},
		g.InitializeScenario,
		nil,
	)
}

func (g *getUsersByIDOrPhone) iHaveAuthenticatedMySelfAsAResourceServer() error {
	g.resourceServer.ID = uuid.New()
	g.resourceServer.Name = "resource_server_test"
	g.resourceServer.Secret = "rs_secret"
	_, err := g.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO resource_servers (id, name, secret) values ('%s', '%s', '%s')", g.resourceServer.ID.String(), g.resourceServer.Name, g.resourceServer.Secret))
	if err != nil {
		return err
	}
	g.apiTest.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(g.resourceServer.ID.String()+":"+g.resourceServer.Secret)))

	return nil
}

func (g *getUsersByIDOrPhone) thereAreUsersWithPhoneNumbers(phones *godog.Table) error {
	phoneData, err := g.apiTest.ReadRowsToMapString(phones)
	if err != nil {
		return err
	}
	for _, v := range phoneData {
		user, err := g.DB.CreateUser(context.Background(), db.CreateUserParams{
			FirstName:  "John",
			MiddleName: "M",
			LastName:   "Doe",
			Phone:      v["phone"],
			UserName:   "jonny",
			Password:   "123456",
			Gender:     "Male",
			ProfilePicture: sql.NullString{
				String: "profile_picture",
				Valid:  true,
			},
			Email: sql.NullString{
				String: fmt.Sprintf("john%s@gmail.com", v),
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
		g.users = append(g.users, user)
	}

	return nil
}

func (g *getUsersByIDOrPhone) iAskForUsersWithPhones(phones *godog.Table) error {
	body, err := g.apiTest.ReadRow(phones, []src.Type{
		{
			Column: "phones",
			Kind:   src.Array,
		},
	}, false)
	if err != nil {
		return err
	}
	g.apiTest.Body = body

	g.apiTest.SendRequest()
	return nil
}

func (g *getUsersByIDOrPhone) iAskForUsersWithIds() error {
	var ids []string
	for _, v := range g.users {
		ids = append(ids, v.ID.String())
	}
	g.apiTest.SetBodyMap(map[string]interface{}{
		"ids": ids,
	})

	g.apiTest.SendRequest()
	return nil
}

func (g *getUsersByIDOrPhone) iShouldGetTheUsers() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	var res dto.RSAPIUsersResponse
	if err := g.apiTest.UnmarshalResponseBodyPath("data", &res); err != nil {
		return err
	}
	ok := true
	for _, v := range res.IDs {
		if !collection.ContainsWithMatcher(v.ID, g.users, func(value uuid.UUID, user db.User) bool {
			return value == user.ID
		}) {
			ok = false
			break
		}
	}
	if ok {
		return nil
	}

	for _, v := range res.Phones {
		if collection.ContainsWithMatcher(v.ID, g.users, func(value uuid.UUID, user db.User) bool {
			return value == user.ID
		}) {
			return fmt.Errorf("expected to find user: %v", v)
		}
	}

	return nil
}

func (g *getUsersByIDOrPhone) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have authenticated my self as a resource server$`, g.iHaveAuthenticatedMySelfAsAResourceServer)
	ctx.Step(`^There are users with phone numbers$`, g.thereAreUsersWithPhoneNumbers)
	ctx.Step(`^I ask for users with phones$`, g.iAskForUsersWithPhones)
	ctx.Step(`^I ask for users with ids$`, g.iAskForUsersWithIds)
	ctx.Step(`^I should get the users$`, g.iShouldGetTheUsers)
	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		g.apiTest.QueryParams = nil
		if _, err := g.Conn.Exec(ctx, fmt.Sprintf("DELETE FROM resource_servers WHERE id='%s'", g.resourceServer.ID.String())); err != nil {
			return ctx, err
		}
		for _, v := range g.users {
			if _, err := g.DB.DeleteUser(ctx, v.ID); err != nil {
				return ctx, err
			}
		}
		g.users = nil

		return ctx, nil
	})
}
