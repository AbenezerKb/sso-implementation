package rest

import "github.com/gin-gonic/gin"

type OAuth interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RequestOTP(ctx *gin.Context)
}

type OAuth2 interface {
	Authorize(ctx *gin.Context)
	GetConsentByID(ctx *gin.Context)
	Approval(ctx *gin.Context)
}
