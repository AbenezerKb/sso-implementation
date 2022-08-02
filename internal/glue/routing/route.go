package routing

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

func RegisterRoutes(group *gin.RouterGroup, routes []Router) {
	for _, route := range routes {
		var handlers []gin.HandlerFunc
		for _, mw := range route.Middlewares {
			handlers = append(handlers, mw)
		}
		handlers = append(handlers, route.Handler)
		group.Handle(route.Method, route.Path, handlers...)
	}
}
