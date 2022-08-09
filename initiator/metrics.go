package initiator

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sso/platform/logger"
)

func InitMetricsRoute(group *gin.Engine, log logger.Logger) {
	group.GET("/metrics", func(context *gin.Context) {
		promhttp.Handler().ServeHTTP(context.Writer, context.Request)
	})
}
