package userinfo

import (
	"context"
	"database/sql"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type userInfoTest struct {
	test.TestInstance
	apiTest src.ApiTest
	user    db.User
}

func TestUserInfo(t *testing.T) {
	u := &userInfoTest{}
	u.TestInstance = test.Initiate("../../../../")
	u.apiTest.URL = "/v1/oauth/userinfo"
	u.apiTest.Method = "GET"
	u.apiTest.SetHeader("Content-Type", "application/json")
	u.apiTest.InitializeServer(u.Server)

	u.apiTest.InitializeTest(t, "User info test", "features/userinfo.feature", u.InitializeScenario)

}

func (u *userInfoTest) thereIsAuthenticatedUserUsingOpenidConnectWithFollowingDetails(userTable *godog.Table) error {
	userJSON, err := u.apiTest.ReadRow(userTable, nil, false)
	if err != nil {
		return err
	}
	var userData dto.User
	err = u.apiTest.UnmarshalJSON([]byte(userJSON), &userData)
	if err != nil {
		return err
	}
	u.user, err = u.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:  userData.FirstName,
		MiddleName: userData.MiddleName,
		LastName:   userData.LastName,
		Email:      sql.NullString{String: userData.Email, Valid: true},
		Phone:      userData.Phone,
		Gender:     userData.Gender,
	})
	if err != nil {
		return err
	}

	accessToken, err := u.PlatformLayer.Token.GenerateAccessToken(context.Background(), u.user.ID.String(), time.Hour)
	if err != nil {
		return err
	}
	u.apiTest.SetHeader("Authorization", "Bearer "+accessToken)

	return nil
}

func (u *userInfoTest) thereIsInvalidAccessToken(accessToken string) error {
	u.apiTest.SetHeader("Authorization", "Bearer "+accessToken)
	return nil
}

func (u *userInfoTest) iSendUserInfoRequest() error {
	u.apiTest.SendRequest()
	return nil
}

func (u *userInfoTest) iShouldGetCorrectUserInfoResponse() error {
	if err := u.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	fetchedUserInfo := dto.UserInfo{}
	err := u.apiTest.UnmarshalResponseBodyPath("data", &fetchedUserInfo)
	if err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedUserInfo.Email, u.user.Email.String); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedUserInfo.FirstName, u.user.FirstName); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedUserInfo.MiddleName, u.user.MiddleName); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedUserInfo.LastName, u.user.LastName); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedUserInfo.Phone, u.user.Phone); err != nil {
		return err
	}

	if err := u.apiTest.AssertEqual(fetchedUserInfo.Gender, u.user.Gender); err != nil {
		return err
	}
	return nil

}
func (u *userInfoTest) theRequestShouldFailWithMessage(message string) error {
	if err := u.apiTest.AssertStatusCode(http.StatusUnauthorized); err != nil {
		return err
	}

	if err := u.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}
	return nil
}

func (u *userInfoTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = u.DB.DeleteUser(context.Background(), u.user.ID)
		return ctx, nil
	})
	ctx.Step(`^I send userInfo request$`, u.iSendUserInfoRequest)
	ctx.Step(`^I should get correct userInfo response$`, u.iShouldGetCorrectUserInfoResponse)
	ctx.Step(`^the request should fail with message "([^"]*)"$`, u.theRequestShouldFailWithMessage)
	ctx.Step(`^there is authenticated user using openid connect with following details$`, u.thereIsAuthenticatedUserUsingOpenidConnectWithFollowingDetails)
	ctx.Step(`^there is invalid access token "([^"]*)"$`, u.thereIsInvalidAccessToken)
}
