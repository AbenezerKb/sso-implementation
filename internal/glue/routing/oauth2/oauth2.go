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
			Method:  "GET",
			Path:    "/approval",
			Handler: handler.ApproveConsent,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:      "GET",
			Path:        "/reject",
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
	}
	routing.RegisterRoutes(oauth2Group, oauth2Routes, enforcer)

}