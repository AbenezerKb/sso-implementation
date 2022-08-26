package oauth2

import (
	"net/http"
	"net/url"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/state"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type oauth2 struct {
	logger       logger.Logger
	oauth2Module module.OAuth2Module
	options      Options
}

type Options struct {
	ErrorURL   string
	ConsentURL string
}

func SetOptions(options Options) Options {
	if options.ErrorURL == "" {
		options.ErrorURL = state.ErrorURL
	}
	if options.ConsentURL == "" {
		options.ConsentURL = state.ConsentURL
	}
	return options
}

func InitOAuth2(logger logger.Logger, oauth2Module module.OAuth2Module, options Options) rest.OAuth2 {
	return &oauth2{
		logger:       logger,
		oauth2Module: oauth2Module,
		options:      options,
	}
}

// Authorize is used to to obtain authorization code.
// @Summary      Authorize.
// @Description  is used to obtain authorization code.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param response_type query string true "response_type"
// @param client_id query string true "client_id"
// @param  state query string true "state"
// @param scope query string true "scope"
// @param redirect_uri query string true "redirect_uri"
// @Success      200
// @Failure      400  {object}  model.ErrorResponse
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/authorize [get]
func (o *oauth2) Authorize(ctx *gin.Context) {
	errorURL, err := url.Parse(o.options.ErrorURL)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "failed to parse error url")
		_ = ctx.Error(err)
		o.logger.Error(ctx, "error parsing error url", zap.Error(err), zap.String("error_url", o.options.ErrorURL))
		return
	}
	errQuery := errorURL.Query()

	authRequestParam := dto.AuthorizationRequestParam{}
	err = ctx.ShouldBindQuery(&authRequestParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "error binding to AuthorizationRequestParam", zap.Error(err), zap.Any("request-uri", ctx.Request.RequestURI))
		errQuery.Set("error", "invalid_request")
		errQuery.Set("error_description", err.Message())
		errQuery.Set("error_code", "400")
		errorURL.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusBadRequest, errorURL.String())
		return
	}

	authRequestParam.ClientID, err = uuid.Parse(ctx.Query("client_id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid client id.")
		o.logger.Info(ctx, "invalid client_id", zap.Error(err), zap.Any("client_id", ctx.Query("client_id")))
		errQuery.Set("error", "invalid_client_id")
		errQuery.Set("error_description", err.Message())
		errQuery.Set("code", "400")

		errorURL.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusFound, errorURL.String())
		return
	}
	requestOrigin := ctx.Request.Host
	if requestOrigin == "" {
		err := errors.ErrInvalidUserInput.New("invalid request origin")
		o.logger.Warn(ctx, "a request without a request origin header was made", zap.Error(err))
		errQuery.Set("error", err.Message())
		errQuery.Set("error_description", err.Error())
		errQuery.Set("error_code", "400")

		errorURL.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusFound, errorURL.String())
		return
	}
	consentId, authErrRsp, err := o.oauth2Module.Authorize(ctx.Request.Context(), authRequestParam, requestOrigin)
	if err != nil {
		o.logger.Info(ctx, "error while authorizing authorization request", zap.Error(err), zap.Any("auth-request-param", authRequestParam))
		errQuery.Set("error", authErrRsp.Error)
		errQuery.Set("error_description", authErrRsp.ErrorDescription)
		errQuery.Set("error_code", "400")

		errorURL.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusFound, errorURL.String())
		return
	}

	consentURL, err := url.Parse(o.options.ConsentURL)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "failed to parse consent url")
		_ = ctx.Error(err)
		o.logger.Error(ctx, "error parsing consent url", zap.Error(err), zap.String("consent_url", o.options.ConsentURL))
		return
	}
	query := consentURL.Query()
	query.Set("consentId", consentId)
	if authRequestParam.Prompt != "" {
		query.Set("prompt", authRequestParam.Prompt)
	} else {
		query.Set("prompt", "consent")
	}

	consentURL.RawQuery = query.Encode()

	o.logger.Info(ctx, "consent url", zap.String("url", consentURL.String()))
	ctx.Redirect(http.StatusFound, consentURL.String())
}

// GetConsentByID is used to get consent by id.
// @Summary      GetConsentByID.
// @Description  is used to get consent by id.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param id path string true "id"
// @param user_id query string true "user_id"
// @Success      200  {object}  dto.ConsentResponse
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /oauth/consent/{id} [get]
// @Security	BearerAuth
func (o *oauth2) GetConsentByID(ctx *gin.Context) {
	consentID := ctx.Param("id")
	consent, err := o.oauth2Module.GetConsentByID(ctx.Request.Context(), consentID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constant.SuccessResponse(ctx, http.StatusOK, consent, nil)
}

// ApproveConsent is used to approve consent.
// @Summary      Consent Approval..
// @Description  is used to approve consent.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param consentId query string true "consentId"
// @success 	 200
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/approveConsent [POST]
// @Security	BearerAuth
func (o *oauth2) ApproveConsent(ctx *gin.Context) {
	consentId := ctx.Query("consentId")
	requestCtx := ctx.Request.Context()
	userIDString, ok := requestCtx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInternalServerError.New("no user_id was found")
		o.logger.Error(ctx, "no user_id was found on gin context", zap.Error(err), zap.String("request-uri", ctx.Request.RequestURI))
		_ = ctx.Error(err)
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "invalid user id")
		o.logger.Error(ctx, "error while parsing x-user-id from request context", zap.Error(err), zap.String("x-user-id", userIDString))
		_ = ctx.Error(err)
		return
	}
	if consentId == "" {
		err := errors.ErrInvalidUserInput.New("invalid consentId")
		o.logger.Info(ctx, "empty consent id", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	opbs, err := ctx.Request.Cookie("opbs")
	if err != nil {
		err := errors.ErrAuthError.Wrap(err, "user not logged in")
		o.logger.Info(ctx, "no opbs value was found while approving authorize request", zap.Error(err))
	}
	redirectURI, err := o.oauth2Module.ApproveConsent(requestCtx, consentId, userID, opbs.Value)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.SetCookie("opbs", utils.GenerateNewOPBS(), 3600, "/", "", true, false)
	ctx.Redirect(http.StatusFound, redirectURI)
}

// RejectConsent is used to reject consent.
// @Summary      Consent Rejection.
// @Description  is used to reject consent.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param consentId query string true "consentId"
// @param failureReason query string true "failureReason"
// @success 	 200
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/rejectConsent [POST]
func (o *oauth2) RejectConsent(ctx *gin.Context) {
	consentId := ctx.Query("consentId")
	failureReason := ctx.GetString("failureReason")
	if consentId == "" {
		err := errors.ErrInvalidUserInput.New("invalid consentId")
		o.logger.Info(ctx, "empty consent id", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	redirectURI, err := o.oauth2Module.RejectConsent(ctx.Request.Context(), consentId, failureReason)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Redirect(http.StatusFound, redirectURI)
}

// Token is used to exchange the authorization code for access token.
// @Summary      exchange token.
// @Description  is used to exchange token.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param tokenParam body dto.AccessTokenRequest true "tokenParam"
// @Success      200  {object}  dto.TokenResponse
// @Failure      404  {object}  model.ErrorResponse "no record of code found"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /oauth/token [post]
// @Security	BasicAuth
func (o *oauth2) Token(ctx *gin.Context) {
	tokenParam := dto.AccessTokenRequest{}
	err := ctx.ShouldBind(&tokenParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Error(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	context := ctx.Request.Context()
	client, ok := context.Value(constant.Context("x-client")).(*dto.Client)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Error(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	resp, err := o.oauth2Module.Token(context, *client, tokenParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constant.SuccessResponse(ctx, http.StatusOK, resp, nil)
}
