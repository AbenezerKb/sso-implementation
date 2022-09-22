package oauth

import (
	"net/http"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.RouterGroup, handler rest.OAuth, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	oauthRoutes := []routing.Router{
		{
			Method:      "POST",
			Path:        "/register",
			Handler:     handler.Register,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
		{
			Method:      "POST",
			Path:        "/login",
			Handler:     handler.Login,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
		{
			Method:      "GET",
			Path:        "/otp",
			Handler:     handler.RequestOTP,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPost,
			Path:    "/logout",
			Handler: handler.Logout,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
			},
			UnAuthorize: true,
		},
		{
			Method:      http.MethodGet,
			Path:        "/refreshToken",
			Handler:     handler.RefreshToken,
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(router, oauthRoutes, enforcer)

}
