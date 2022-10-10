package get_scope

import (
	"context"
	"database/sql"
	"net/http"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/test"
	"testing"

	"github.com/cucumber/godog"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type getScopeTest struct {
	test.TestInstance
	apiTest src.ApiTest
	scopes  []db.Scope
	admin   db.User
}

func TestGetScope(t *testing.T) {
	g := &getScopeTest{}
	g.apiTest.Method = http.MethodGet
	g.apiTest.SetHeader("Content-Type", "application/json")
	g.TestInstance = test.Initiate("../../../../")
	g.apiTest.InitializeServer(g.Server)
	g.apiTest.InitializeTest(t, "get scope test", "features/get_scope.feature", g.InitializeScenario)
}

func (g *getScopeTest) iAmLoggedInAsAdminUser(adminCredentials *godog.Table) error {
	var err error
	g.admin, err = g.Authenticate(adminCredentials)
	if err != nil {
		return err
	}
	_, g.GrantRoleAfterFunc, err = g.GrantRoleForUserWithAfter(g.admin.ID.String(), adminCredentials)
	if err != nil {
		return err
	}

	g.apiTest.SetHeader("Authorization", "Bearer "+g.AccessToken)

	return nil
}
func (g *getScopeTest) theFollowingScopesAreRegisteredOnTheSystem(scopes *godog.Table) error {
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

func (g *getScopeTest) iRequestToGetAScopeBy(name string) error {
	g.apiTest.URL += name
	g.apiTest.SendRequest()
	return nil
}

func (g *getScopeTest) iShouldGetTheFollowingScope(scope *godog.Table) error {
	var responseScope dto.Scope

	if err := g.apiTest.AssertStatusCode(http.StatusOK); err != nil {
		return err
	}

	scopeJSON, err := g.apiTest.ReadRow(scope, nil, false)
	if err != nil {
		return err
	}

	var scopeData dto.Scope
	err = g.apiTest.UnmarshalJSON([]byte(scopeJSON), &scopeData)
	if err != nil {
		return err
	}

	err = g.apiTest.UnmarshalResponseBodyPath("data", &responseScope)
	if err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(responseScope.Name, scopeData.Name); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(responseScope.Description, scopeData.Description); err != nil {
		return err
	}

	if err := g.apiTest.AssertEqual(responseScope.ResourceServerName, scopeData.ResourceServerName); err != nil {
		return err
	}

	return nil
}

func (g *getScopeTest) myRequestShouldFailWith(message string) error {
	if err := g.apiTest.AssertStatusCode(http.StatusNotFound); err != nil {
		return err
	}
	return g.apiTest.AssertStringValueOnPathInResponse("error.message", message)
}

func (g *getScopeTest) InitializeScenario(ctx *godog.ScenarioContext) {
	g.apiTest.URL = "/v1/oauth/scopes/"

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		_, _ = g.DB.DeleteUser(ctx, g.admin.ID)

		for _, v := range g.scopes {
			_, _ = g.DB.DeleteScope(ctx, v.Name)
		}
		_ = g.GrantRoleAfterFunc()

		return ctx, nil
	})

	ctx.Step(`^I am logged in as admin user$`, g.iAmLoggedInAsAdminUser)
	ctx.Step(`^I request to get a scope by "([^"]*)"$`, g.iRequestToGetAScopeBy)
	ctx.Step(`^I should get the following scope$`, g.iShouldGetTheFollowingScope)
	ctx.Step(`^my request should fail with "([^"]*)"$`, g.myRequestShouldFailWith)
	ctx.Step(`^The following scopes are registered on the system$`, g.theFollowingScopesAreRegisteredOnTheSystem)
}
