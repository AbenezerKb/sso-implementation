package oauth2

import (
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
		},
		{
			Method:      "GET",
			Path:        "/consent/:id",
			Handler:     handler.GetConsentByID,
			Middlewares: []gin.HandlerFunc{},
		},

		{
			Method:      "GET",
			Path:        "/consent",
			Handler:     handler.Approval,
			Middlewares: []gin.HandlerFunc{},
		},
	}
	routing.RegisterRoutes(oauth2Group, oauth2Routes, enforcer)

}
