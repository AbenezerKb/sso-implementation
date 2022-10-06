package middleware

import (
	"context"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/permissions"
	"sso/internal/module"
	"sso/platform"
	"sso/platform/logger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type AuthMiddleware interface {
	Authentication() gin.HandlerFunc
	AccessControl() gin.HandlerFunc
	ClientBasicAuth() gin.HandlerFunc
	MiniRideBasicAuth() gin.HandlerFunc
}

type MiniRideCredential struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type authMiddleware struct {
	enforcer           *casbin.Enforcer
	auth               module.OAuthModule
	token              platform.Token
	client             module.ClientModule
	miniRideCredential MiniRideCredential
	role               module.RoleModule
	logger             logger.Logger
}

func InitAuthMiddleware(enforcer *casbin.Enforcer,
	auth module.OAuthModule, token platform.Token, client module.ClientModule, miniRideCredential MiniRideCredential, role module.RoleModule, logger logger.Logger) AuthMiddleware {
	return &authMiddleware{
		enforcer,
		auth,
		token,
		client,
		miniRideCredential,
		role,
		logger,
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
		valid, claims := a.token.VerifyToken(jwt.SigningMethodPS512, tokenString)
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
		requestCtx := ctx.Request.Context()
		userId := requestCtx.Value(constant.Context("x-user-id")).(string)

		// check user role status
		status, err := a.role.GetRoleStatusForUser(ctx, userId)
		if err != nil {
			err := errors.ErrAcessError.Wrap(err, "unable to perform operation")
			_ = ctx.Error(err)
			a.logger.Error(ctx, "error while fetching role status for user", zap.Error(err), zap.String("user-id", userId))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if status != constant.Active {
			message := "access denied"
			if status != "" {
				message = "your role is disabled"
			}
			err := errors.ErrAcessError.New(message)
			a.logger.Info(ctx, "access denied", zap.Error(err), zap.String("user-id", userId))
			_ = ctx.Error(err)
			ctx.AbortWithStatus(http.StatusForbidden)
		}

		// enforce role
		err = a.enforcer.LoadPolicy()
		if err != nil {
			err := errors.ErrAcessError.Wrap(err, "unable to perform operation")
			_ = ctx.Error(err)
			a.logger.Error(ctx, "error while trying to load policy", zap.Error(err), zap.String("user-id", userId))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ok, err := a.enforcer.Enforce(userId, permissions.Notneeded, permissions.Notneeded, ctx.Request.URL.Path, ctx.Request.Method, permissions.Notneeded)
		if err != nil {
			err := errors.ErrAcessError.Wrap(err, "unable to perform operation")
			_ = ctx.Error(err)
			a.logger.Error(ctx, "error while enforcing policy", zap.Error(err), zap.String("user-id", userId))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !ok {
			err := errors.ErrAcessError.Wrap(err, "Access denied")
			_ = ctx.Error(err)
			a.logger.Info(ctx, "access denied", zap.Error(err), zap.String("user-id", userId))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}

func (a *authMiddleware) ClientBasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		clientId, secret, ok := ctx.Request.BasicAuth()
		if !ok {
			err := errors.ErrInternalServerError.New("could not get extract client credentials")
			a.logger.Error(ctx, "extract error", zap.Error(err))
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		client, err := a.client.GetClientByID(ctx.Request.Context(), clientId)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}

		if client.Status != constant.Active {
			Err := errors.ErrAuthError.Wrap(nil, "Your account has been deactivated, Please activate your account.")
			ctx.Error(Err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ok := client.Secret == secret; !ok {
			err = errors.ErrAcessError.Wrap(err, "unauthorized_client")
			a.logger.Info(ctx, "unauthorized_client", zap.Error(err), zap.String("client-secret", client.Secret), zap.String("provided-secret", secret))
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), constant.Context("x-client"), client))
		ctx.Next()
	}
}

func (a *authMiddleware) MiniRideBasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, password, ok := ctx.Request.BasicAuth()
		if !ok {
			err := errors.ErrInternalServerError.New("couldn't extract basic auth detail's")
			a.logger.Error(ctx, "extract error", zap.Error(err))
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if password != a.miniRideCredential.Password || username != a.miniRideCredential.UserName {
			err := errors.ErrAuthError.New("mini_ride credential mismatch")
			a.logger.Info(ctx, "mini_ride_credential_mismatch", zap.Error(err), zap.Any("existing_credential", a.miniRideCredential), zap.Any("provided_credential", []string{username, password}))
			ctx.Error(err)
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
