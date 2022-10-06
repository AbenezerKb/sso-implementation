package profile

import (
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.RouterGroup, handler rest.Profile, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	profile := router.Group("/profile")
	profileRoutes := []routing.Router{
		{
			Method:  "PUT",
			Path:    "",
			Handler: handler.UpdateProfile,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  "GET",
			Path:    "",
			Handler: handler.GetProfile,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  "PUT",
			Path:    "/picture",
			Handler: handler.UpdateProfilePicture,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  "PATCH",
			Path:    "/phone",
			Handler: handler.ChangePhone,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(profile, profileRoutes, enforcer)
}
