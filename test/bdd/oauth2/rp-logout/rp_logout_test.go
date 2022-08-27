package logout

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
	"sso/test"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type rpLogoutTest struct {
	test.TestInstance
	apiTest  src.ApiTest
	user     db.User
	client   db.Client
	id_token string
}

type logoutRspQueryParams struct {
	PostLogoutRedirectUri string `json:"post_logout_redirect_uri"`
	State                 string `json:"state"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
}

func TestLogout(t *testing.T) {
	r := &rpLogoutTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest.InitializeTest(t, "Logout test", "features/rp_logout.feature", r.InitializeScenario)
}

func (r *rpLogoutTest) iAmRegisteredOnTheSystem() error {

	var err error
	if r.client, err = r.DB.CreateClient(context.Background(), db.CreateClientParams{
		RedirectUris: utils.ArrayToString([]string{"https://www.google.com"}),
		Name:         "google",
		Scopes:       "openid",
		ClientType:   "confidential",
		Secret:       utils.GenerateRandomString(25, true),
		LogoUrl:      "https://www.google.com/images/errors/robot.png",
	}); err != nil {
		return err
	}
	return nil
}

func (r *rpLogoutTest) iHaveId_token() error {
	var err error
	r.id_token, err = r.PlatformLayer.Token.GenerateIdToken(context.Background(), &dto.User{
		ID:         r.user.ID,
		FirstName:  r.user.FirstName,
		Email:      r.user.Email.String,
		MiddleName: r.user.MiddleName,
	}, r.client.ID.String(), time.Hour*24)
	if err != nil {
		return err
	}
	return nil
}

func (r *rpLogoutTest) iHaveTheFollowingDetails(logoutParams *godog.Table) error {
	state, _ := r.apiTest.ReadCellString(logoutParams, "state")
	logout_redirect_uri, _ := r.apiTest.ReadCellString(logoutParams, "post_logout_redirect_uri")

	r.apiTest.SetQueryParam("state", state)
	r.apiTest.SetQueryParam("post_logout_redirect_uri", logout_redirect_uri)
	r.apiTest.SetQueryParam("id_token_hint", r.id_token)
	return nil
}

func (r *rpLogoutTest) iHaveTheFollowingInvalid_requestDetails(logoutParams *godog.Table) error {
	state, _ := r.apiTest.ReadCellString(logoutParams, "state")
	logout_redirect_uri, _ := r.apiTest.ReadCellString(logoutParams, "post_logout_redirect_uri")
	id_token, _ := r.apiTest.ReadCellString(logoutParams, "id_token_hint")

	r.apiTest.SetQueryParam("state", state)
	r.apiTest.SetQueryParam("post_logout_redirect_uri", logout_redirect_uri)
	r.apiTest.SetQueryParam("id_token_hint", id_token)
	return nil
}

func (r *rpLogoutTest) iRequestToLogout() error {
	r.apiTest.SetHeader("Authorization", "Basic "+basicAuth(r.client.ID.String(), r.client.Secret))

	r.apiTest.SendRequest()
	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (r *rpLogoutTest) iShouldBeRedirectedToWithTheFollowingQueryParams(logout_uri string, logoutParams *godog.Table) error {
	if err := r.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}

	location := r.apiTest.Response.Header().Get("Location")
	parsedLocation, err := url.Parse(location)
	if err != nil {
		return err
	}

	rawPath := fmt.Sprintf("%s://%s%s", parsedLocation.Scheme, parsedLocation.Host, parsedLocation.Path)
	if err := r.apiTest.AssertEqual(rawPath, logout_uri); err != nil {
		return err
	}

	query := parsedLocation.Query()
	if query.Has("post_logout_redirect_uri") != true {
		return fmt.Errorf("expected post_logout_redirect_uri in post_logout_redirect_uri query parameter")
	}
	return nil
}

func (r *rpLogoutTest) theUserIsRegisteredOnTheSystem() error {
	var err error
	hash, err := utils.HashAndSalt(context.Background(), []byte("password"), r.Logger)
	if err != nil {
		return err
	}
	if r.user, err = r.DB.CreateUser(context.Background(), db.CreateUserParams{
		Email:      utils.StringOrNull("yonaskemon@gmail.com"),
		Password:   hash,
		FirstName:  "someone",
		MiddleName: "someone",
		LastName:   "someone",
		Phone:      "0987654321",
	}); err != nil {
		return err
	}

	return nil
}

func (r *rpLogoutTest) iShouldBeRedirectedToWithTheFollowingFailureQueryParams(err_uri string, errParams *godog.Table) error {
	if err := r.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}

	param, err := r.apiTest.ReadRow(errParams, nil, false)
	if err != nil {
		return err
	}
	var rspParamsQuery logoutRspQueryParams
	err = r.apiTest.UnmarshalJSONAt([]byte(param), "", &rspParamsQuery)
	if err != nil {
		return err
	}

	location := r.apiTest.Response.Header().Get("Location")
	parsedLocation, err := url.Parse(location)
	if err != nil {
		return err
	}

	rawPath := fmt.Sprintf("%s://%s%s", parsedLocation.Scheme, parsedLocation.Host, parsedLocation.Path)
	if err := r.apiTest.AssertEqual(rawPath, err_uri); err != nil {
		return err
	}

	query := parsedLocation.Query()
	if query.Has("error") != true {
		return fmt.Errorf("expected error in error query parameter")
	}
	if query.Has("error_description") != true {
		return fmt.Errorf("expected error_description in error_description query parameter")
	}

	if err := r.apiTest.AssertEqual(query.Get("error"), rspParamsQuery.Error); err != nil {
		return err
	}

	if err := r.apiTest.AssertEqual(query.Get("error_description"), rspParamsQuery.ErrorDescription); err != nil {
		return err
	}

	return nil
}

func (r *rpLogoutTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		r.apiTest.URL = "/v1/oauth/logout"
		r.apiTest.Method = http.MethodGet
		r.apiTest.SetHeader("Content-Type", "application/json")
		r.apiTest.InitializeServer(r.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = r.DB.DeleteUser(ctx, r.user.ID)
		return ctx, err
	})

	ctx.Step(`^I am  registered on the system$`, r.iAmRegisteredOnTheSystem)
	ctx.Step(`^I have id_token$`, r.iHaveId_token)
	ctx.Step(`^I have the following details:$`, r.iHaveTheFollowingDetails)
	ctx.Step(`^I have the following invalid_request details:$`, r.iHaveTheFollowingInvalid_requestDetails)
	ctx.Step(`^I request to logout$`, r.iRequestToLogout)
	ctx.Step(`^I should be redirected to "([^"]*)" with the following failure query params:$`, r.iShouldBeRedirectedToWithTheFollowingFailureQueryParams)
	ctx.Step(`^I should be redirected to "([^"]*)" with the following query params:$`, r.iShouldBeRedirectedToWithTheFollowingQueryParams)
	ctx.Step(`^the user is registered on the system$`, r.theUserIsRegisteredOnTheSystem)
}
