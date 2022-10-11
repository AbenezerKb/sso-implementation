package create_identity_provider

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
	"net/http"
	"reflect"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"
)

type createIdentityProviderTest struct {
	test.TestInstance
	apiTest          src.ApiTest
	identityProvider dto.IdentityProvider
	Admin            db.User
}

func TestCreateIdentityProvider(t *testing.T) {
	r := &createIdentityProviderTest{}
	r.TestInstance = test.Initiate("../../../../")
	r.apiTest = src.ApiTest{
		Server: r.Server,
	}
	r.apiTest.URL = "/v1/identityProviders"
	r.apiTest.Method = http.MethodPost
	r.apiTest.SetHeader("Content-Type", "application/json")
	r.apiTest.InitializeTest(t, "create identity provider test", "features/create_identity_provider.feature", r.InitializeScenario)
}

// background functions
func (i *createIdentityProviderTest) iAmLoggedInWithTheFollowingCredentials(adminCredentials *godog.Table) error {
	var err error
	i.Admin, err = i.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, i.GrantRoleAfterFunc, err = i.GrantRoleForUserWithAfter(i.Admin.ID.String(), adminCredentials)
	return err
}

// given functions
func (i *createIdentityProviderTest) iHaveFilledTheFollowingDataForTheIdentityProvider(ipTable *godog.Table) error {
	ipJSON, err := i.apiTest.ReadRow(ipTable, nil, false)
	if err != nil {
		return err
	}
	i.apiTest.Body = ipJSON
	return i.apiTest.UnmarshalJSON([]byte(ipJSON), &i.identityProvider)
}

// when functions
func (i *createIdentityProviderTest) iSubmitToCreateAnIdentityProvider() error {
	i.apiTest.SetHeader("Authorization", "Bearer "+i.AccessToken)
	i.apiTest.SendRequest()

	return nil
}

// then functions
func (i *createIdentityProviderTest) theIdentityProviderShouldBeCreated() error {
	if err := i.apiTest.AssertStatusCode(http.StatusCreated); err != nil {
		return err
	}
	var identityProvider dto.IdentityProvider
	if err := i.apiTest.UnmarshalResponseBodyPath("data", &identityProvider); err != nil {
		return err
	}
	i.identityProvider.ID = identityProvider.ID
	i.identityProvider.CreatedAt = identityProvider.CreatedAt
	i.identityProvider.UpdatedAt = identityProvider.UpdatedAt
	i.identityProvider.Status = identityProvider.Status
	if !reflect.DeepEqual(identityProvider, i.identityProvider) {
		return fmt.Errorf("got %v, \nwant %v", identityProvider, i.identityProvider)
	}
	return nil
}
func (i *createIdentityProviderTest) theRequestShouldFailWith(message string) error {
	if err := i.apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
		return err
	}
	if err := i.apiTest.AssertStringValueOnPathInResponse("error.field_error.0.description", message); err != nil {
		return err
	}
	return nil
}

func (i *createIdentityProviderTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = i.DB.DeleteUser(ctx, i.Admin.ID)
		_, _ = i.DB.DeleteIdentityProvider(ctx, i.identityProvider.ID)
		_ = i.GrantRoleAfterFunc()
		return ctx, nil
	})
	ctx.Step(`^I am logged in with the following credentials$`, i.iAmLoggedInWithTheFollowingCredentials)
	ctx.Step(`^I have filled the following data for the identity provider$`, i.iHaveFilledTheFollowingDataForTheIdentityProvider)
	ctx.Step(`^I submit to create an identity provider$`, i.iSubmitToCreateAnIdentityProvider)
	ctx.Step(`^the identity provider should be created$`, i.theIdentityProviderShouldBeCreated)
	ctx.Step(`^the request should fail with "([^"]*)"$`, i.theRequestShouldFailWith)
}
