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
		o.logger.Error(ctx, zap.Error(err).String)
		ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	registeredUser, err := o.oauthModule.Register(ctx.Request.Context(), userParam)
	if err != nil {
		ctx.Error(err)
		return
	}
	o.logger.Info(ctx, "registered user")
	constant.SuccessResponse(ctx, http.StatusCreated, registeredUser, nil)
}

func (o *oauth) Login(ctx *gin.Context) {
	userParam := dto.User{}
	err := ctx.ShouldBind(&userParam)

	if err != nil {
		o.logger.Error(ctx, "invalid input", zap.Error(err))
		ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}

	loginRsp, err := o.oauthModule.Login(ctx.Request.Context(), userParam)

	if err != nil {
		ctx.Error(err)
		return
	}

	//Todo: save session
	//Todo: possibly redirect to authorize endpoint

	ctx.SetCookie("access_token", loginRsp.AccessToken, 3600, "/", "", false, true)
	ctx.SetCookie("id_token", loginRsp.RefreshToken, 12000, "/", "", false, true)
	o.logger.Info(ctx, "user logged in")

	redirectUrl := ctx.Query("redirect_url")
	if redirectUrl == "" {
		redirectUrl = "/"
	}
	ctx.Redirect(http.StatusFound, redirectUrl)
	// constant.SuccessResponse(ctx, http.StatusOK, loginRsp, nil)
}
