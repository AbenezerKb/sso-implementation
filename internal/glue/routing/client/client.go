package client

import (
	"net/http"
	"sso/internal/constant/permissions"
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, client rest.Client, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	clients := group.Group("/clients")
	clientRoutes := []routing.Router{
		{
			Method:  "POST",
			Path:    "",
			Handler: client.CreateClient,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.CreateClient,
		},
		{
			Method:  "DELETE",
			Path:    "/:id",
			Handler: client.DeleteClient,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.DeleteClient,
		},
		{
			Method:  http.MethodGet,
			Path:    "",
			Handler: client.GetAllClients,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetAllClients,
		},
		{
			Method:  http.MethodGet,
			Path:    "/:id",
			Handler: client.GetAllClientByID,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.GetClient,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/:id/status",
			Handler: client.UpdateClientStatus,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.UpdateClient,
		},
		{
			Method:  http.MethodPut,
			Path:    "/:id",
			Handler: client.UpdateClient,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.Authentication(),
				authMiddleware.AccessControl(),
			},
			Permission: permissions.UpdateClient,
		},
	}
	routing.RegisterRoutes(clients, clientRoutes, enforcer)
}
