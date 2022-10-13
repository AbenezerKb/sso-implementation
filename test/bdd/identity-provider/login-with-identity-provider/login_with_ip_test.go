package login_with_identity_provider

import (
	"context"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/mocks/platform/identityProvider"
	"sso/platform"
	"sso/test"
	"testing"
)

type loginWithIPTest struct {
	test.TestInstance
	admin   db.User
	apiTest src.ApiTest
	ip      db.IdentityProvider
	ipUser  struct {
		ID string `json:"id"`
		dto.User
	}
	ipServer    platform.IdentityProvider
	requestCode string
}

func TestLoginWithIP(t *testing.T) {
	l := loginWithIPTest{}
	l.TestInstance = test.Initiate("../../../../")
	l.apiTest = src.ApiTest{
		Server: l.Server,
		URL:    "/v1/loginWithIP",
		Method: http.MethodPost,
	}
	l.apiTest.SetHeader("Content-Type", "application/json")
	l.apiTest.InitializeTest(t, "login with ip test", "features/login_with_ip.feature", l.InitializeScenario)
}

// background
func (l *loginWithIPTest) thereExistsAnIdentityProviderWithTheFollowingInfo(providerTable *godog.Table) error {
	// save identity provider to database
	providerJSON, err := l.apiTest.ReadRow(providerTable, nil, false)
	if err != nil {
		return err
	}
	var providerData dto.IdentityProvider
	err = l.apiTest.UnmarshalJSON([]byte(providerJSON), &providerData)
	if err != nil {
		return err
	}
	l.ip, err = l.DB.CreateIdentityProvider(context.Background(), db.CreateIdentityProviderParams{
		Name:             providerData.Name,
		ClientID:         providerData.ClientID,
		ClientSecret:     providerData.ClientSecret,
		TokenEndpointUrl: providerData.TokenEndpointURI,
	})
	if err != nil {
		return err
	}

	return nil
}
func (l *loginWithIPTest) iAmRegisteredOnThatIdentityProviderAsFollows(userTable *godog.Table) error {
	// register user for the identity provider
	userJSON, err := l.apiTest.ReadRow(userTable, nil, false)
	if err != nil {
		return err
	}

	err = l.apiTest.UnmarshalJSON([]byte(userJSON), &l.ipUser)
	if err != nil {
		return err
	}
	return err
}

// given
func (l *loginWithIPTest) iHaveGrantedConsentToMyLoginWithCode(code string) error {
	// register code for the identity provider
	l.ipServer = identityProvider.InitIP(l.ip.ClientID, l.ip.ClientSecret, "veryLegitCode", "legit-access-token", dto.UserInfo{
		Sub:            l.ipUser.ID,
		FirstName:      l.ipUser.FirstName,
		MiddleName:     l.ipUser.MiddleName,
		LastName:       l.ipUser.LastName,
		Email:          l.ipUser.Email,
		Phone:          l.ipUser.Phone,
		Gender:         l.ipUser.Gender,
		ProfilePicture: l.ipUser.ProfilePicture,
	})
	if err := identityProvider.SetUserForProvider(dto.UserInfo{
		Sub:            l.ipUser.ID,
		FirstName:      l.ipUser.FirstName,
		MiddleName:     l.ipUser.MiddleName,
		LastName:       l.ipUser.LastName,
		Email:          l.ipUser.Email,
		Phone:          l.ipUser.Phone,
		Gender:         l.ipUser.Gender,
		ProfilePicture: l.ipUser.ProfilePicture,
	}, &l.PlatformLayer.SelfIP); err != nil {
		return err
	}
	l.requestCode = code
	return nil
}

// when
func (l *loginWithIPTest) iRequestToLoginWithIdentityProvider(provider string) error {
	if provider == l.ip.Name {
		provider = l.ip.ID.String()
	}
	l.apiTest.SetBodyMap(map[string]interface{}{
		"code": l.requestCode,
		"ip":   provider,
	})
	l.apiTest.SendRequest()

	return nil
}

// then
func (l *loginWithIPTest) iShouldSuccessfullyLogin() error {
	if err := l.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	if err := l.apiTest.AssertColumnExists("data.refresh_token"); err != nil {
		return err
	}
	if err := l.apiTest.AssertColumnExists("data.access_token"); err != nil {
		return err
	}
	if err := l.apiTest.AssertColumnExists("data.id_token"); err != nil {
		return err
	}
	return nil
}
func (l *loginWithIPTest) myRequestShouldFailWithMessage(message string) error {
	if err := l.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}

	if err := l.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}

	return nil
}
func (l *loginWithIPTest) myRequestShouldFailWith(message string) error {
	status := http.StatusBadRequest
	if message == "authentication failed" {
		status = http.StatusUnauthorized
	}
	if err := l.apiTest.AssertStatusCode(status); err != nil {
		return err
	}

	if err := l.apiTest.AssertStringValueOnPathInResponse("error.message", message); err != nil {
		return err
	}

	return nil
}

func (l *loginWithIPTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = l.DB.DeleteUser(ctx, l.admin.ID)
		_, _ = l.Conn.Exec(ctx, "DELETE FROM users WHERE phone = $1", l.ipUser.Phone)
		_, _ = l.DB.DeleteIdentityProvider(ctx, l.ip.ID)
		return ctx, nil
	})
	ctx.Step(`^I am registered on that identity provider as follows$`, l.iAmRegisteredOnThatIdentityProviderAsFollows)
	ctx.Step(`^I should successfully login$`, l.iShouldSuccessfullyLogin)
	ctx.Step(`^There exists an identity provider with the following info$`, l.thereExistsAnIdentityProviderWithTheFollowingInfo)
	ctx.Step(`^I have granted consent to my login with code "([^"]*)"$`, l.iHaveGrantedConsentToMyLoginWithCode)
	ctx.Step(`^I request to login with identity provider "([^"]*)"$`, l.iRequestToLoginWithIdentityProvider)
	ctx.Step(`^my request should fail with "([^"]*)"$`, l.myRequestShouldFailWith)
	ctx.Step(`^my request should fail with message "([^"]*)"$`, l.myRequestShouldFailWithMessage)
}
