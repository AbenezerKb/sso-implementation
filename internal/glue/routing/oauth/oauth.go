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
			Method:      http.MethodPost,
			Path:        "/register",
			Handler:     handler.Register,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
		{
			Method:      http.MethodPost,
			Path:        "/login",
			Handler:     handler.Login,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
		{
			Method:      http.MethodGet,
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
		{
			Method:      http.MethodPost,
			Path:        "/loginWithIP",
			Handler:     handler.LoginWithIP,
			UnAuthorize: true,
		},
		{
			Method:      http.MethodGet,
			Path:        "/registeredIdentityProviders",
			Handler:     handler.GetIdentityProviders,
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(router, oauthRoutes, enforcer)

}
