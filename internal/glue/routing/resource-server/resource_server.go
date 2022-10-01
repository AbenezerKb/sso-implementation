package resource_server

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"
)

func InitRoute(group *gin.RouterGroup, resourceServer rest.ResourceServer, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	resourceServers := group.Group("/resourceServers")
	resourceServerRoutes := []routing.Router{
		{
			Method:  http.MethodPost,
			Path:    "",
			Handler: resourceServer.CreateResourceServer,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.CreateResourceServer,
		},
		{
			Method:  http.MethodGet,
			Path:    "",
			Handler: resourceServer.GetAllResourceServers,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetAllResourceServers,
		},
	}

	routing.RegisterRoutes(resourceServers, resourceServerRoutes, enforcer)
}
