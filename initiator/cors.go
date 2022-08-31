package initiator

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func InitCORS() gin.HandlerFunc {
	origins := viper.GetStringSlice("cors.origin")
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	allowCredentials := viper.GetString("cors.allow_credentials")
	if allowCredentials == "" {
		allowCredentials = "true"
	}
	headers := viper.GetStringSlice("cors.headers")
	if len(headers) == 0 {
		headers = []string{"*"}
	}
	methods := viper.GetStringSlice("cors.methods")
	if len(methods) == 0 {
		methods = []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		}
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", allowCredentials)
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
