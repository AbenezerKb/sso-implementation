package get_user_by_phone

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
	"sso/test"
	"testing"
)

type getUserByIDOrPhone struct {
	test.TestInstance
	apiTest        src.ApiTest
	user           db.User
	resourceServer db.ResourceServer
}

func TestGetUserByIDOrPhone(t *testing.T) {
	g := getUserByIDOrPhone{}
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.Server = g.Server
	g.apiTest.URL = "/v1/internal/user"
	g.apiTest.RunTest(t,
		"get user by phone test",
		&src.TestOptions{
			Paths: []string{
				"features/get_user_by_phone.feature",
				"features/get_user_by_id.feature",
			},
		},
		g.InitializeScenario,
		nil,
	)
}

// given
func (g *getUserByIDOrPhone) iHaveAuthenticatedMySelfAsAResourceServer() error {
	g.resourceServer.ID = uuid.New()
	g.resourceServer.Name = "resource_server_test"
	g.resourceServer.Secret = "rs_secret"
	_, err := g.Conn.Exec(context.Background(), fmt.Sprintf("INSERT INTO resource_servers (id, name, secret) values ('%s', '%s', '%s')", g.resourceServer.ID.String(), g.resourceServer.Name, g.resourceServer.Secret))
	if err != nil {
		return err
	}

	return nil
}

func (g *getUserByIDOrPhone) thereIsAUserWithPhoneNumber(phone string) error {
	user, err := g.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:  "John",
		MiddleName: "M",
		LastName:   "Doe",
		Phone:      phone,
		UserName:   "jonny",
		Password:   "123456",
		Gender:     "Male",
		ProfilePicture: sql.NullString{
			String: "profile_picture",
			Valid:  true,
		},
		Email: sql.NullString{
			String: "john@gmail.com",
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	g.user = user

	return nil
}

// when
func (g *getUserByIDOrPhone) iAskForAUserWithPhoneNumber(phone string) error {
	g.apiTest.SetQueryParam("phone", phone)
	g.apiTest.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(g.resourceServer.ID.String()+":"+g.resourceServer.Secret)))
	g.apiTest.SendRequest()

	return nil
}

func (g *getUserByIDOrPhone) iAskForAUserWithId() error {
	g.apiTest.SetQueryParam("id", g.user.ID.String())
	g.apiTest.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(g.resourceServer.ID.String()+":"+g.resourceServer.Secret)))
	g.apiTest.SendRequest()

	return nil
}

func (g *getUserByIDOrPhone) iAskForAUserWithIncorrectId() error {
	g.apiTest.SetQueryParam("id", uuid.NewString())
	g.apiTest.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(g.resourceServer.ID.String()+":"+g.resourceServer.Secret)))
	g.apiTest.SendRequest()

	return nil
}

// then
func (g *getUserByIDOrPhone) iShouldGetTheUserData() error {
	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	var user dto.User
	if err := g.apiTest.UnmarshalResponseBodyPath("data", &user); err != nil {
		return err
	}

	if user.ID != g.user.ID || user.FirstName != g.user.FirstName || user.MiddleName != g.user.MiddleName ||
		user.LastName != g.user.LastName || user.Phone != g.user.Phone || user.Email != g.user.Email.String ||
		user.Status != g.user.Status.String || user.ProfilePicture != g.user.ProfilePicture.String {
		return fmt.Errorf("got %v, want %v", user, g.user)
	}

	return nil
}

func (g *getUserByIDOrPhone) myRequestShouldFailWithMessage(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}

	if err := g.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (g *getUserByIDOrPhone) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I ask for a user with incorrect id$`, g.iAskForAUserWithIncorrectId)
	ctx.Step(`^I ask for a user with id$`, g.iAskForAUserWithId)
	ctx.Step(`^I ask for a user with phone number "([^"]*)"$`, g.iAskForAUserWithPhoneNumber)
	ctx.Step(`^I have authenticated my self as a resource server$`, g.iHaveAuthenticatedMySelfAsAResourceServer)
	ctx.Step(`^I should get the user data$`, g.iShouldGetTheUserData)
	ctx.Step(`^My request should fail with message "([^"]*)"$`, g.myRequestShouldFailWithMessage)
	ctx.Step(`^There is a user with phone number "([^"]*)"$`, g.thereIsAUserWithPhoneNumber)
	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		g.apiTest.QueryParams = nil
		if _, err := g.Conn.Exec(ctx, fmt.Sprintf("DELETE FROM resource_servers WHERE id='%s'", g.resourceServer.ID.String())); err != nil {
			return ctx, err
		}
		if _, err := g.DB.DeleteUser(ctx, g.user.ID); err != nil {
			return ctx, err
		}

		return ctx, nil
	})
}
