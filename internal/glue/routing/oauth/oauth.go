package oauth

import (
	"sso/internal/glue/routing"
	"sso/internal/handler/rest"

	"github.com/gin-gonic/gin"
)

func InitRoute(router *gin.RouterGroup, handler rest.OAuth) {
	// router.POST("/register", handler.Register)
	oauthRoutes := []routing.Router{
		{
			Method:      "POST",
			Path:        "/register",
			Handler:     handler.Register,
			Middlewares: []gin.HandlerFunc{},
		},
	}
	routing.RegisterRoutes(router, oauthRoutes)

}
