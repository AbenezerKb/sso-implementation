package oauth

import (
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
	"sso/platform/utils"

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

// Register creates a new user.
// @Summary      Register a new user.
// @Description  Register a new user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @param user body dto.RegisterUser true "user"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /register [post]
func (o *oauth) Register(ctx *gin.Context) {
	userParam := dto.RegisterUser{}
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
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

// Login logs in a user.
// @Summary      Login a user.
// @Description  Login a user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @param login_credential body dto.LoginCredential true "login_credential"
// @Success      200  {object}  dto.TokenResponse
// @Failure      401  {object}  model.ErrorResponse "invalid credentials"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /login [post]
func (o *oauth) Login(ctx *gin.Context) {
	userParam := dto.LoginCredential{}
	err := ctx.ShouldBind(&userParam)

	if err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}

	loginRsp, err := o.oauthModule.Login(ctx.Request.Context(), userParam)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.SetCookie("opbs", utils.GenerateNewOPBS(), 3600, "/", "", true, false)
	ctx.SetCookie("access_token", loginRsp.AccessToken, 3600, "/", "", false, true)
	ctx.SetCookie("refresh_token", loginRsp.RefreshToken, 12000, "/", "", false, true)
	o.logger.Info(ctx, "user logged in")

	constant.SuccessResponse(ctx, http.StatusOK, loginRsp, nil)
}

// RequestOTP is used to request otp.
// @Summary      Request otp.
// @Description  is used to request otp for login and signup
// @Tags         auth
// @Accept       json
// @Produce      json
// @param phone query string true "phone"
// @param type query string true "type can be login or signup" Enums(login, signup)
// @Success      200  {boolean}  true
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /otp [get]
func (o *oauth) RequestOTP(ctx *gin.Context) {
	phone := ctx.Query("phone")
	RqType := ctx.Query("type")
	if phone == "" || RqType == "" {
		o.logger.Info(ctx, "invalid input", zap.String("phone", phone))
		_ = ctx.Error(errors.ErrInvalidUserInput.New("invalid phone"))
		return
	}
	err := o.oauthModule.RequestOTP(ctx.Request.Context(), phone, RqType)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	o.logger.Info(ctx, "OTP sent", zap.String("phone", phone))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// Logout logs out a user.
// @Summary      logout  user.
// @Description  logout user.
// @Tags         auth
// @param tokenParam body dto.InternalRefreshTokenRequestBody true "logoutParam"
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /logout [post]
// @Security	BearerAuth
func (o *oauth) Logout(ctx *gin.Context) {
	refreshTokenRequest := dto.InternalRefreshTokenRequestBody{}
	if err := ctx.ShouldBind(&refreshTokenRequest); err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	err := o.oauthModule.Logout(ctx.Request.Context(), refreshTokenRequest)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	// change opbs
	ctx.SetCookie("opbs", utils.GenerateNewOPBS(), 3600, "/", "", true, false)
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// RefreshToken refreshs a user access token.
// @Summary      refresh access token.
// @Description  refresh access token.
// @Tags         auth
// @param tokenParam body dto.InternalRefreshTokenRequestBody true "refreshTokenParam"
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /refreshtoken [post]
func (o *oauth) RefreshToken(ctx *gin.Context) {
	refreshTokenRequest := dto.InternalRefreshTokenRequestBody{}
	if err := ctx.ShouldBind(&refreshTokenRequest); err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}

	resp, err := o.oauthModule.RefreshToken(ctx.Request.Context(), refreshTokenRequest)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constant.SuccessResponse(ctx, http.StatusOK, resp, nil)
}
