package role

import (
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
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
		{
			Method:  http.MethodGet,
			Path:    "",
			Handler: handler.GetAllRoles,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetAllRoles,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/:name/status",
			Handler: handler.UpdateRoleStatus,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.ChangeRoleStatus,
		},
		{
			Method:  http.MethodGet,
			Path:    "/:name",
			Handler: handler.GetRoleByName,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetRole,
		},
	}
	routing.RegisterRoutes(roleGroup, roleRoutes, enforcer)
}
