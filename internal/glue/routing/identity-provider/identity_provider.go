package identity_provider

import (
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, identityProvider rest.IdentityProvider, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	identityProviders := group.Group("/identityProviders")
	identityProviderRoutes := []routing.Router{
		{
			Method:  http.MethodPost,
			Path:    "",
			Handler: identityProvider.CreateIdentityProvider,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.CreateIdentityProvider,
		},
		{
			Method:  http.MethodPut,
			Path:    "/:id",
			Handler: identityProvider.UpdateIdentityProvider,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.UpdateIdentityProvider,
		},
		{
			Method:  http.MethodGet,
			Path:    "/:id",
			Handler: identityProvider.GetIdentityProvider,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetIdentityProvider,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/:id",
			Handler: identityProvider.DeleteIdentityProvider,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.DeleteIdentityProvider,
		},
	}

	routing.RegisterRoutes(identityProviders, identityProviderRoutes, enforcer)
}
