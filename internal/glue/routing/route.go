package routing

import (
	"path"
	"sso/internal/constant/permissions"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
	Permission  string
}

func RegisterRoutes(group *gin.RouterGroup, routes []Router, enforcer *casbin.Enforcer) {
	for _, route := range routes {
		var handler []gin.HandlerFunc
		handler = append(handler, route.Middlewares...)
		handler = append(handler, route.Handler)
		group.Handle(route.Method, route.Path, handler...)

		if len(route.Middlewares) > 0 {
			url := path.Join(group.BasePath(), route.Path)

			if exists := enforcer.HasPolicy("PERMISSION", route.Permission, permissions.PermissionCategory[route.Permission], url, route.Method); !exists {
				enforcer.AddPolicy("PERMISSION", route.Permission, permissions.PermissionCategory[route.Permission], url, route.Method)
			}
		}
	}
}
