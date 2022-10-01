package mini_ride

import (
	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.RouterGroup, handler rest.MiniRide, authMiddleware middleware.AuthMiddleware, enforcer *casbin.Enforcer) {
	miniRide := router.Group("/users")
	miniRideRoutes := []routing.Router{
		{
			Method:  "GET",
			Path:    "/exists/:phone",
			Handler: handler.CheckPhone,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.MiniRideBasicAuth(),
			},
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(miniRide, miniRideRoutes, enforcer)
}
