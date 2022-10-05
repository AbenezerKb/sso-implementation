package rest

import "github.com/gin-gonic/gin"

type OAuth interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RequestOTP(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type OAuth2 interface {
	Authorize(ctx *gin.Context)
	GetConsentByID(ctx *gin.Context)
	ApproveConsent(ctx *gin.Context)
	RejectConsent(ctx *gin.Context)
	Token(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RevokeClient(ctx *gin.Context)
	GetAuthorizedClients(ctx *gin.Context)
	GetOpenIDAuthorizedClients(ctx *gin.Context)
}
type User interface {
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	UpdateUserStatus(ctx *gin.Context)
}

type Client interface {
	CreateClient(ctx *gin.Context)
	DeleteClient(ctx *gin.Context)
	GetAllClients(ctx *gin.Context)
	GetAllClientByID(ctx *gin.Context)
	UpdateClientStatus(ctx *gin.Context)
	UpdateClient(ctx *gin.Context)
}

type Scope interface {
	GetScope(ctx *gin.Context)
	CreateScope(ctx *gin.Context)
	GetAllScopes(ctx *gin.Context)
	DeleteScope(ctx *gin.Context)
}

type Profile interface {
	UpdateProfile(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	UpdateProfilePicture(ctx *gin.Context)
}

type MiniRide interface {
	CheckPhone(ctx *gin.Context)
}
type ResourceServer interface {
	CreateResourceServer(ctx *gin.Context)
	GetAllResourceServers(ctx *gin.Context)
}

type Role interface {
	GetAllPermissions(ctx *gin.Context)
	CreateRole(ctx *gin.Context)
}
