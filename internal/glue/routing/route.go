package routing

import (
	"path"
	"sso/internal/constant"
	"sso/internal/constant/permissions"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
	Permission  permissions.Permission
	UnAuthorize bool
}

func RegisterRoutes(group *gin.RouterGroup, routes []Router, enforcer *casbin.Enforcer) {
	for _, route := range routes {
		var handler []gin.HandlerFunc
		handler = append(handler, route.Middlewares...)
		handler = append(handler, route.Handler)
		group.Handle(route.Method, route.Path, handler...)

		if !route.UnAuthorize {
			url := path.Join(group.BasePath(), route.Path)

			if exists := enforcer.HasPolicy(route.Permission.ID); !exists {
				enforcer.AddPolicy(route.Permission.ID, route.Permission.Name, route.Permission.Category, url, route.Method, constant.Active)
			}
		}
	}
}
