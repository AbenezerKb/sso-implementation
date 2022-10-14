package user

import (
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.RouterGroup, handler rest.User, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	users := router.Group("/users")
	userRoutes := []routing.Router{
		{
			Method:  "POST",
			Path:    "",
			Handler: handler.CreateUser,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.CreateUser,
		},
		{
			Method:  "GET",
			Path:    "/:id",
			Handler: handler.GetUser,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "",
			Handler: handler.GetAllUsers,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetAllUsers,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/:id/status",
			Handler: handler.UpdateUserStatus,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.UpdateUserStatus,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/:id/role",
			Handler: handler.UpdateUserRole,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.UpdateUserRole,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/:id/role",
			Handler: handler.RevokeUserRole,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.RevokeUserRole,
		},
	}
	routing.RegisterRoutes(users, userRoutes, enforcer)
}
