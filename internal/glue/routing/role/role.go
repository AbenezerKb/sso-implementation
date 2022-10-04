package role

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"
)

func InitRoute(group *gin.RouterGroup, handler rest.Role, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	roleGroup := group.Group("roles")
	roleRoutes := []routing.Router{
		{
			Method:  "GET",
			Path:    "permissions",
			Handler: handler.GetAllPermissions,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetAllPermissions,
		},
		{
			Method:  http.MethodPost,
			Path:    "",
			Handler: handler.CreateRole,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.CreateRole,
		},
	}
	routing.RegisterRoutes(roleGroup, roleRoutes, enforcer)
}
