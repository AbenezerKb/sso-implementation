package identity_provider

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"
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
	}

	routing.RegisterRoutes(identityProviders, identityProviderRoutes, enforcer)
}
