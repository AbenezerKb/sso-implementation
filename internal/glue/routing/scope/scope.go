package scope

import (
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, handler rest.Scope, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	scopeGroup := group.Group("oauth/scope")
	scopeRoutes := []routing.Router{
		{
			Method:  "GET",
			Path:    "/:name",
			Handler: handler.GetScope,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			UnAuthorize: true,
		},
		{
			Method:  "POST",
			Path:    "",
			Handler: handler.CreateScope,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(scopeGroup, scopeRoutes, enforcer)
}
