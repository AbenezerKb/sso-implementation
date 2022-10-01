package get_scopes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getScopesTest struct {
	test.TestInstance
	apiTest     src.ApiTest
	scopes      []db.Scope
	Admin       db.User
	Preferences preferenceData
}

type preferenceData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func TestGetScopes(t *testing.T) {
	g := getScopesTest{}
	g.apiTest.URL = "/v1/scopes"
	g.apiTest.Method = http.MethodGet
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get scopes test", "features/get_scopes.feature", g.InitializeScenario)
}

func (g *getScopesTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	g.Admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	return g.GrantRoleForUser(g.Admin.ID.String(), adminCredentials)
}

func (g *getScopesTest) theFollowingScopesAreRegisteredOnTheSystem(scopes *godog.Table) error {
	scopesData, err := g.apiTest.ReadRows(scopes, nil, false)
	if err != nil {
		return err
	}
	var scopesStruct []dto.Scope
	if err := g.apiTest.UnmarshalJSONAt([]byte(scopesData), "", &scopesStruct); err != nil {
		return err
	}
	for _, scope := range scopesStruct {
		savedScope, err := g.DB.CreateScope(context.Background(), db.CreateScopeParams{
			Name:        scope.Name,
			Description: scope.Description,
			ResourceServerName: sql.NullString{
				String: scope.ResourceServerName,
				Valid:  true,
			},
		})
		if err != nil {
			return err
		}
		g.scopes = append(g.scopes, savedScope)
	}
	return nil
}

func (g *getScopesTest) iRequestToGetAllTheScopesWithTheFollowingPreferences(preferences *godog.Table) error {
	preferencesJSON, err := g.apiTest.ReadRow(preferences, []src.Type{
		{
			Column: "page",
			Kind:   src.Any,
		},
		{
			Column: "per_page",
			Kind:   src.Any,
		},
	}, false)
	if err != nil {
		return err
	}
	err = g.apiTest.UnmarshalJSON([]byte(preferencesJSON), &g.Preferences)
	if err != nil {
		return err
	}

	g.apiTest.SetQueryParam("page", fmt.Sprintf("%d", g.Preferences.Page))
	g.apiTest.SetQueryParam("per_page", fmt.Sprintf("%d", g.Preferences.PerPage))
	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)
	g.apiTest.SendRequest()
	return nil
}

func (g *getScopesTest) iShouldGetTheListOfScopesThatPassMyPreferences() error {
	var responseScopes []dto.Scope
	var metaData model.MetaData

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}
	err := g.apiTest.UnmarshalResponseBodyPath("meta_data", &metaData)
	if err != nil {
		return err
	}

	err = g.apiTest.UnmarshalResponseBodyPath("data", &responseScopes)
	if err != nil {
		return err
	}
	var total int
	if g.Preferences.Page < metaData.Total/g.Preferences.PerPage {
		total = g.Preferences.PerPage
	} else if g.Preferences.Page == metaData.Total/g.Preferences.PerPage {
		total = metaData.Total % g.Preferences.PerPage
	} else {
		total = 0
	}
	if err := g.apiTest.AssertEqual(len(responseScopes), total); err != nil {
		return err
	}
	for _, v := range responseScopes {
		found := false
		for _, v2 := range g.scopes {
			if v.Name == v2.Name {
				found = true
				continue
			}
		}
		if !found {
			return fmt.Errorf("expected client: %v", v)
		}
	}
	return nil
}

func (g *getScopesTest) InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get all the scopes with the following preferences$`, g.iRequestToGetAllTheScopesWithTheFollowingPreferences)
	ctx.Step(`^I should get the list of scopes that pass my preferences$`, g.iShouldGetTheListOfScopesThatPassMyPreferences)
	ctx.Step(`^The following scopes are registered on the system$`, g.theFollowingScopesAreRegisteredOnTheSystem)
}
