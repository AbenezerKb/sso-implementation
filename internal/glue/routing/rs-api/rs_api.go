package rs_api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"
)

func InitRoute(router *gin.RouterGroup, handler rest.RSAPI, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	internal := router.Group("/internal")
	internalRoutes := []routing.Router{
		{
			Method:  http.MethodGet,
			Path:    "/user",
			Handler: handler.GetUserByPhoneOrID,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.ResourceServerBasicAuth(),
			},
			UnAuthorize: true,
		},
		{
			Method:  http.MethodPost,
			Path:    "users",
			Handler: handler.GetUsersByPhoneOrID,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.ResourceServerBasicAuth(),
			},
			UnAuthorize: true,
		},
	}

	routing.RegisterRoutes(internal, internalRoutes, enforcer)
}
