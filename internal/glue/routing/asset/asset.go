package asset

import (
	"net/http"

	"sso/internal/glue/routing"
	"sso/internal/handler/middleware"
	"sso/internal/handler/rest"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func InitRoute(
	group *gin.RouterGroup,
	asset rest.Asset,
	_ middleware.AuthMiddleware,
	enforcer *casbin.Enforcer,
) {
	assetGroup := group.Group("assets")
	assetGroup.Static("", "assets")
	assetRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Handler:     asset.UploadAsset,
			Middlewares: []gin.HandlerFunc{},
			UnAuthorize: true,
		},
	}

	routing.RegisterRoutes(assetGroup, assetRoutes, enforcer)
}
