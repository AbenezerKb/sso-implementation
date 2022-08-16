package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type authorizationTest struct {
	test.TestInstance
	apiTest      src.ApiTest
	requestParam *dto.AuthorizationRequestParam
}

type authRspQueryParams struct {
	consentId        string `json:"consentId"`
	State            string `json:"state"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func TestAuthorization(t *testing.T) {
	a := &authorizationTest{}
	a.TestInstance = test.Initiate("../../../")
	a.apiTest.InitializeTest(t, "Authorization test", "features/authorization.feature", a.InitializeScenario)

}

func (a *authorizationTest) iHaveTheFollowingParameters(params *godog.Table) error {
	param, err := a.apiTest.ReadRow(params, nil, false)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(param), &a.requestParam)
	if err != nil {
		return err
	}

	a.apiTest.SetQueryParam("client_id", a.requestParam.ClientID.String())
	a.apiTest.SetQueryParam("response_type", a.requestParam.ResponseType)
	a.apiTest.SetQueryParam("state", a.requestParam.State)
	a.apiTest.SetQueryParam("scope", a.requestParam.Scope)
	a.apiTest.SetQueryParam("redirect_uri", a.requestParam.RedirectURI)

	return nil
}

func (a *authorizationTest) iSendAPOSTRequest() error {
	a.apiTest.SendRequest()
	return nil
}

func (a *authorizationTest) iShouldBeRedirectedToWithTheFollowingSuccessParameters(redirect_uri string, rspParams *godog.Table) error {
	if err := a.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}
	param, err := a.apiTest.ReadRow(rspParams, nil, false)
	if err != nil {
		return err
	}
	var rspParamsQuery authRspQueryParams
	err = a.apiTest.UnmarshalJSONAt([]byte(param), "", &rspParamsQuery)
	if err != nil {
		return err
	}

	location := a.apiTest.Response.Header().Get("Location")
	parsedLocation, err := url.Parse(location)
	if err != nil {
		return err
	}

	rawPath := fmt.Sprintf("%s://%s%s", parsedLocation.Scheme, parsedLocation.Host, parsedLocation.Path)
	if err := a.apiTest.AssertEqual(rawPath, redirect_uri); err != nil {
		return err
	}

	query := parsedLocation.Query()
	if query.Has("consentId") != true {
		return fmt.Errorf("expected consentId in consentId")
	}

	return nil
}

func (a *authorizationTest) iShouldBeRedirectedToWithTheFollowingErrorParameters(redirect_uri string, rspParams *godog.Table) error {
	if err := a.apiTest.AssertStatusCode(http.StatusFound); err != nil {
		return err
	}
	param, err := a.apiTest.ReadRow(rspParams, nil, false)
	if err != nil {
		return err
	}
	var rspParamsQuery authRspQueryParams
	err = a.apiTest.UnmarshalJSONAt([]byte(param), "", &rspParamsQuery)
	if err != nil {
		return err
	}

	location := a.apiTest.Response.Header().Get("Location")
	parsedLocation, err := url.Parse(location)
	if err != nil {
		return err
	}

	query := parsedLocation.Query()
	if err := a.apiTest.AssertEqual(query.Get("error"), rspParamsQuery.Error); err != nil {
		return err
	}
	if err := a.apiTest.AssertEqual(query.Get("error_description"), rspParamsQuery.ErrorDescription); err != nil {
		return err
	}

	return nil
}

func (a *authorizationTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		a.apiTest.URL = "/v1/oauth/authorize"
		a.apiTest.Method = "GET"
		a.apiTest.SetHeader("Content-Type", "application/json")
		a.apiTest.InitializeServer(a.Server)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})

	ctx.Step(`^I have the following parameters:$`, a.iHaveTheFollowingParameters)
	ctx.Step(`^I send a POST request$`, a.iSendAPOSTRequest)
	ctx.Step(`^I should be redirected to "([^"]*)" with the following error parameters:$`, a.iShouldBeRedirectedToWithTheFollowingErrorParameters)
	ctx.Step(`^I should be redirected to "([^"]*)" with the following success parameters:$`, a.iShouldBeRedirectedToWithTheFollowingSuccessParameters)

}
