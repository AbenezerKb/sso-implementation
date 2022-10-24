package profile

import (
	"net/http"
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
			Method:  http.MethodPut,
			Path:    "",
			Handler: handler.UpdateProfile,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodGet,
			Path:    "",
			Handler: handler.GetProfile,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPut,
			Path:    "/picture",
			Handler: handler.UpdateProfilePicture,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/phone",
			Handler: handler.ChangePhone,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/password",
			Handler: handler.ChangePassword,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/devices",
			Handler: handler.GetAllCurrentSessions,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(profile, profileRoutes, enforcer)
}
