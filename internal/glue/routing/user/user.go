package user

import (
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
			Method:  "PATCH",
			Path:    "",
			Handler: handler.UpdateProfile,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
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
	}
	routing.RegisterRoutes(users, userRoutes, enforcer)
}
