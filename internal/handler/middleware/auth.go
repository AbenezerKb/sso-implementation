package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/permissions"
	"sso/internal/module"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware interface {
	Authentication() gin.HandlerFunc
	AccessControl() gin.HandlerFunc
}

type authMiddleware struct {
	enforcer     *casbin.Enforcer
	auth         module.OAuthModule
	ssoPublicKey *rsa.PublicKey
}

func InitAuthMiddleware(enforcer *casbin.Enforcer,
	auth module.OAuthModule, ssoPublicKey *rsa.PublicKey) AuthMiddleware {
	return &authMiddleware{
		enforcer,
		auth,
		ssoPublicKey,
	}
}

func (a *authMiddleware) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			Err := errors.ErrInvalidToken.New("Unauthorized")
			ctx.Error(Err)
			ctx.Abort()
			return
		}

		tokenString := authHeader[len(bearer):]
		valid, claims := a.auth.VerifyToken(jwt.SigningMethodPS512, tokenString, a.ssoPublicKey)
		if !valid {
			Err := errors.ErrAuthError.New("Unauthorized")
			ctx.Error(Err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userStatus, err := a.auth.GetUserStatus(ctx.Request.Context(), claims.Subject)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		if userStatus != constant.Active {
			Err := errors.ErrAuthError.Wrap(nil, "Your account has been deactivated, Please activate your account.")
			ctx.Error(Err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), constant.Context("x-user-id"), claims.Subject))
		ctx.Next()
	}
}
func (a *authMiddleware) AccessControl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context := ctx.Request.Context()
		userId := context.Value(constant.Context("x-user-id")).(string)

		a.enforcer.LoadPolicy()
		ok, err := a.enforcer.Enforce(userId, permissions.Notneeded, permissions.Notneeded, ctx.Request.URL.Path, ctx.Request.Method)
		if err != nil {
			Err := errors.ErrAcessError.Wrap(err, "unable to perform operation")
			ctx.Error(Err)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !ok {
			Err := errors.ErrAcessError.Wrap(err, "Access denied")
			ctx.Error(Err)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
