package oauth2

import (
	"net/http"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, handler rest.OAuth2, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	oauth2Group := group.Group("/oauth")
	oauth2Routes := []routing.Router{
		{
			Method:      "GET",
			Path:        "/authorize",
			Handler:     handler.Authorize,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
		{
			Method:  "GET",
			Path:    "/consent/:id",
			Handler: handler.GetConsentByID,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},

		{
			Method:  "POST",
			Path:    "/approveConsent",
			Handler: handler.ApproveConsent,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:      "POST",
			Path:        "/rejectConsent",
			Handler:     handler.RejectConsent,
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/token",
			Handler: handler.Token,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.ClientBasicAuth(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/logout",
			Handler: handler.Logout,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.ClientBasicAuth(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/revokeClient",
			Handler: handler.RevokeClient,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/authorizedClients",
			Handler: handler.GetAuthorizedClients,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/openIDAuthorizedClients",
			Handler: handler.GetOpenIDAuthorizedClients,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/userinfo",
			Handler: handler.UserInfo,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(oauth2Group, oauth2Routes, enforcer)

}
