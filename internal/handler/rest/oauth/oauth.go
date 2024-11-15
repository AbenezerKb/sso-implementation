package oauth

import (
	"net/http"

	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
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
	options     Options
}

type Options struct {
	RefreshTokenCookie utils.CookieOptions
	OPBSCookie         utils.CookieOptions
}

func SetOptions(options Options) Options {
	if options.OPBSCookie.Path == "" {
		options.OPBSCookie.Path = "/"
	}
	if options.OPBSCookie.MaxAge == 0 {
		options.OPBSCookie.MaxAge = 365 * 24 * 60 * 60
	}
	if options.OPBSCookie.SameSite < 1 || options.OPBSCookie.SameSite > 4 {
		options.OPBSCookie.SameSite = 4
	}

	if options.RefreshTokenCookie.Path == "" {
		options.RefreshTokenCookie.Path = "/"
	}
	if options.RefreshTokenCookie.MaxAge == 0 {
		options.RefreshTokenCookie.MaxAge = 365 * 24 * 60 * 60
	}
	if options.RefreshTokenCookie.SameSite < 1 || options.RefreshTokenCookie.SameSite > 4 {
		options.RefreshTokenCookie.SameSite = 3
	}

	return options
}
func InitOAuth(logger logger.Logger, oauthModule module.OAuthModule, options Options) rest.OAuth {
	return &oauth{
		logger:      logger,
		oauthModule: oauthModule,
		options:     options,
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

	loginRsp, err := o.oauthModule.Login(ctx.Request.Context(), userParam, dto.UserDeviceAddress{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	})

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	utils.SetOPBSCookie(ctx, utils.GenerateNewOPBS(), o.options.OPBSCookie)
	utils.SetRefreshTokenCookie(ctx, loginRsp.RefreshToken, o.options.RefreshTokenCookie)
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

// RequestResetCode is used to request reset code.
// @Summary      Request reset code.
// @Description  is used to request reset code for forget password
// @Tags         auth
// @Accept       json
// @Produce      json
// @param email query string true "email"
// @Success      200  {boolean}  true
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /resetCode [get]
func (o *oauth) RequestResetCode(ctx *gin.Context) {
	err := o.oauthModule.RequestResetCode(ctx.Request.Context(), ctx.Query("email"))
	if err != nil {
		_ = ctx.Error(err)

		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// ResetPassword is used to reset password
// @Summary      reset password.
// @Description  is used to reset password for forgotten password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @param request body dto.ResetPasswordRequest true "request"
// @Success      200  {boolean}  true
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /resetPassword [post]
func (o *oauth) ResetPassword(ctx *gin.Context) {
	var resetPasswordRequest dto.ResetPasswordRequest

	err := ctx.ShouldBind(&resetPasswordRequest)
	if err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))

		return
	}

	err = o.oauthModule.ResetPassword(ctx, resetPasswordRequest)
	if err != nil {
		_ = ctx.Error(err)

		return
	}

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
	utils.SetOPBSCookie(ctx, utils.GenerateNewOPBS(), o.options.OPBSCookie)
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
// @Router       /refreshToken [get]
func (o *oauth) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("ab_fen")
	if err != nil {
		o.logger.Info(ctx, "no refresh token was found", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "no refresh token was found"))
		return
	}

	resp, err := o.oauthModule.RefreshToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		_ = ctx.Error(err)
		utils.RemoveRefreshTokenCookie(ctx, o.options.RefreshTokenCookie)
		return
	}

	utils.SetRefreshTokenCookie(ctx, resp.RefreshToken, o.options.RefreshTokenCookie)
	constant.SuccessResponse(ctx, http.StatusOK, resp, nil)
}

// LoginWithIP logs in a user with an identity provider.
// @Summary      Login a user with an identity provider.
// @Description  Login a user with an identity provider.
// @Tags         auth
// @Accept       json
// @Produce      json
// @param login_with_ip body request_models.LoginWithIP true "login_with_ip"
// @Success      200  {object}  dto.TokenResponse
// @Failure      401  {object}  model.ErrorResponse "invalid credentials"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /loginWithIP [post]
func (o *oauth) LoginWithIP(ctx *gin.Context) {
	var login request_models.LoginWithIP
	err := ctx.ShouldBind(&login)
	if err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}

	loginRsp, err := o.oauthModule.LoginWithIdentityProvider(ctx.Request.Context(), login, dto.UserDeviceAddress{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	})

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	utils.SetOPBSCookie(ctx, utils.GenerateNewOPBS(), o.options.OPBSCookie)
	utils.SetRefreshTokenCookie(ctx, loginRsp.RefreshToken, o.options.RefreshTokenCookie)
	o.logger.Info(ctx, "user logged in")

	constant.SuccessResponse(ctx, http.StatusOK, loginRsp, nil)
}

// GetIdentityProviders fetches all identity provider that user can login.
// @Summary      get all identity providers.
// @Description  get all identity providers.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  []dto.IdentityProvider
// @Failure 	400 {object} model.ErrorResponse
// @Router 		/registeredIdentityProviders [get]
func (o *oauth) GetIdentityProviders(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	idPs, err := o.oauthModule.GetAllIdentityProviders(requestCtx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, idPs, nil)
}
