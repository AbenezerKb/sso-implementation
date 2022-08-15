package rest

import "github.com/gin-gonic/gin"

type OAuth interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RequestOTP(ctx *gin.Context)
}

type User interface {
	CreateUser(ctx *gin.Context)
}

type Client interface {
	CreateClient(ctx *gin.Context)
}
