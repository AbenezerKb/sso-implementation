package oauth2

import (
	"net/http"
	"net/url"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
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
}

func InitOAuth2(logger logger.Logger, oauth2Module module.OAuth2Module) rest.OAuth2 {
	return &oauth2{
		logger:       logger,
		oauth2Module: oauth2Module,
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
	requestCtx := ctx.Request.Context()

	authRequestParam := dto.AuthorizationRequestParam{}
	err := ctx.ShouldBindQuery(&authRequestParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid_request")
		o.logger.Info(ctx, "error binding to AuthorizationRequestParam", zap.Error(err), zap.Any("request-uri", ctx.Request.RequestURI))

		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.Authorize(requestCtx, authRequestParam, "", err))
		return
	}

	authRequestParam.ClientID, err = uuid.Parse(ctx.Query("client_id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid_client_id")
		o.logger.Info(ctx, "invalid client_id", zap.Error(err), zap.Any("client_id", ctx.Query("client_id")))

		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.Authorize(requestCtx, authRequestParam, "", err))
		return
	}

	requestOrigin := ctx.Request.Host
	if requestOrigin == "" {
		err := errors.ErrInvalidUserInput.New("invalid request origin")
		o.logger.Warn(ctx, "a request without a request origin header was made", zap.Error(err))

		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.Authorize(requestCtx, authRequestParam, requestOrigin, err))
		return
	}

	ctx.Redirect(
		http.StatusFound,
		o.oauth2Module.Authorize(requestCtx, authRequestParam, requestOrigin, nil))
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
// @param consent_id body string true "consent_id"
// @success 	 200
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/approveConsent [POST]
// @Security	BearerAuth
func (o *oauth2) ApproveConsent(ctx *gin.Context) {
	var consentResultRsp = dto.ConsentResultRsp{}
	err := ctx.ShouldBind(&consentResultRsp)
	if err != nil {
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}

	requestCtx := ctx.Request.Context()
	userIDString, ok := requestCtx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInternalServerError.New("no user_id was found")
		o.logger.Error(ctx, "no user_id was found on gin context", zap.Error(err), zap.String("request-uri", ctx.Request.RequestURI))
		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, uuid.UUID{}, "", err))
		return
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "invalid user id")
		o.logger.Error(ctx, "error while parsing x-user-id from request context", zap.Error(err), zap.String("x-user-id", userIDString))
		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, uuid.UUID{}, "", err))
		return
	}
	if consentResultRsp.ConsentID == "" {
		err := errors.ErrInvalidUserInput.New("invalid consentId")
		o.logger.Info(ctx, "empty consent id", zap.Error(err))
		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, userID, "", err))
		return
	}

	opbs, err := ctx.Request.Cookie("opbs")
	if err != nil {
		err := errors.ErrAuthError.Wrap(err, "user not logged in")
		o.logger.Warn(ctx, "no opbs value was found while approving authorize request", zap.Error(err))
		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, userID, "", err))
		return
	}

	ctx.SetCookie("opbs", utils.GenerateNewOPBS(), 3600, "/", "", true, false)
	ctx.Redirect(
		http.StatusFound,
		o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, userID, opbs.Value, nil))
}

// RejectConsent is used to reject consent.
// @Summary      Consent Rejection.
// @Description  is used to reject consent.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param consent_id body string true "consent_id"
// @param failureReason query string true "failureReason"
// @success 	 200
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/rejectConsent [POST]
func (o *oauth2) RejectConsent(ctx *gin.Context) {
	consentId := ctx.Query("consentId")
	failureReason := ctx.Query("failureReason")
	if consentId == "" {
		err := errors.ErrInvalidUserInput.New("invalid consentId")
		o.logger.Info(ctx, "empty consent id", zap.Error(err))
		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.RejectConsent(ctx, consentId, "", err))
		return
	}

	ctx.Redirect(
		http.StatusFound,
		o.oauth2Module.RejectConsent(ctx.Request.Context(), consentId, failureReason, nil))
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
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	context := ctx.Request.Context()
	client, ok := context.Value(constant.Context("x-client")).(*dto.Client)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
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

// Logout is used to logout user.
// @Summary      rp-logout.
// @Description  this is requested from truest client only.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param id_token_hint query string true "id_token_hint"
// @param post_logout_redirect_uri query string true "post_logout_redirect_uri"
// @param state query string true "state"
// @Success      200
// @Failure      400  {object}  model.ErrorResponse
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/logout [get]
// @Security	BasicAuth
func (o *oauth2) Logout(ctx *gin.Context) {
	errRedirectUri, err := url.Parse("error url")
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "invalid error uri")
		o.logger.Error(ctx, "invalid error uri", zap.Error(err))
		_ = ctx.Error(err)
	}
	errQuery := errRedirectUri.Query()

	logoutReqParam := dto.LogoutRequest{}
	err = ctx.ShouldBindQuery(&logoutReqParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))

		errQuery.Set("error", "invalid request")
		errQuery.Set("error_description", "no logedin user found")
		errRedirectUri.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusFound, errRedirectUri.String())

		return
	}

	redirectURI, errRsp, err := o.oauth2Module.Logout(ctx.Request.Context(), logoutReqParam)
	if err != nil {
		errQuery.Set("error", errRsp.Error)
		errQuery.Set("error_description", errRsp.ErrorDescription)
		errRedirectUri.RawQuery = errQuery.Encode()

		ctx.Redirect(http.StatusFound, errRedirectUri.String())
		return
	}

	ctx.SetCookie("opbs", utils.GenerateNewOPBS(), 3600, "/", "", true, false)
	ctx.Redirect(http.StatusFound, redirectURI)
}
