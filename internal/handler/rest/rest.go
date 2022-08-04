package rest

import "github.com/gin-gonic/gin"

type OAuth interface {
	Register(ctx *gin.Context)
}
