package oauth

import (
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type oauth struct {
	logger      logger.Logger
	oauthModule module.OAuthModule
}

func InitOAuth(logger logger.Logger, oauthModule module.OAuthModule) rest.OAuth {
	return &oauth{
		logger,
		oauthModule,
	}
}

// implement Oauth
func (o *oauth) Register(ctx *gin.Context) {
	userParam := dto.User{}
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		o.logger.Info(ctx, zap.Error(err).String)
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	registeredUser, err := o.oauthModule.Register(ctx.Request.Context(), userParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	o.logger.Info(ctx, "registered user")
	constant.SuccessResponse(ctx, http.StatusCreated, registeredUser, nil)
}
