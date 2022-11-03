package oauth2

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
	"strings"

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
	OPBSCookie utils.CookieOptions
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

	requestOrigin := strings.TrimSuffix(ctx.Request.Header.Get("Referer"), "/")
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
// @success 	 200 {object} dto.RedirectResponse "redirect response"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/approveConsent [POST]
// @Security	BearerAuth
func (o *oauth2) ApproveConsent(ctx *gin.Context) {
	var consentResultRsp = dto.ConsentResultRsp{}
	requestCtx := ctx.Request.Context()

	err := ctx.ShouldBind(&consentResultRsp)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, uuid.UUID{}, "", err),
			}, nil)
		return
	}

	userIDString, ok := requestCtx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInternalServerError.New("no user_id was found")
		o.logger.Error(ctx, "no user_id was found on gin context", zap.Error(err), zap.String("request-uri", ctx.Request.RequestURI))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, uuid.UUID{}, "", err),
			}, nil)
		return
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "invalid user id")
		o.logger.Error(ctx, "error while parsing x-user-id from request context", zap.Error(err), zap.String("x-user-id", userIDString))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, uuid.UUID{}, "", err),
			}, nil)
		return
	}
	if consentResultRsp.ConsentID == "" {
		err := errors.ErrInvalidUserInput.New("invalid consentId")
		o.logger.Info(ctx, "empty consent id", zap.Error(err))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, userID, "", err),
			}, nil)
		return
	}

	opbs, err := ctx.Request.Cookie("opbs")
	if err != nil {
		err := errors.ErrAuthError.Wrap(err, "user not logged in")
		o.logger.Warn(ctx, "no opbs value was found while approving authorize request", zap.Error(err))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, userID, "", err),
			}, nil)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK,
		dto.RedirectResponse{
			Location: o.oauth2Module.ApproveConsent(requestCtx, consentResultRsp.ConsentID, userID, opbs.Value, nil),
		}, nil)
}

// RejectConsent is used to reject consent.
// @Summary      Consent Rejection.
// @Description  is used to reject consent.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param consent_id body string true "consent_id"
// @param failure_reason body string true "failure_reason"
// @success 	 200 {object} dto.RedirectResponse "redirect response"
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/rejectConsent [POST]
func (o *oauth2) RejectConsent(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()

	var consentResultRsp = dto.ConsentResultRsp{}
	err := ctx.ShouldBind(&consentResultRsp)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.RejectConsent(requestCtx, consentResultRsp.ConsentID, "", err),
			}, nil)
		return
	}
	if consentResultRsp.ConsentID == "" {
		err := errors.ErrInvalidUserInput.New("invalid consentId")
		o.logger.Info(ctx, "empty consent id", zap.Error(err))
		constant.SuccessResponse(ctx, http.StatusOK,
			dto.RedirectResponse{
				Location: o.oauth2Module.RejectConsent(requestCtx, consentResultRsp.ConsentID, "", err),
			}, nil)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK,
		dto.RedirectResponse{
			Location: o.oauth2Module.RejectConsent(requestCtx, consentResultRsp.ConsentID, consentResultRsp.FailureReason, nil),
		}, nil)
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

	logoutReqParam := dto.LogoutRequest{}
	requestCtx := ctx.Request.Context()
	err := ctx.ShouldBindQuery(&logoutReqParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid request")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		ctx.Redirect(
			http.StatusFound,
			o.oauth2Module.Logout(requestCtx, logoutReqParam, err))
		return
	}

	utils.SetOPBSCookie(ctx, utils.GenerateNewOPBS(), o.options.OPBSCookie)
	ctx.Redirect(
		http.StatusFound,
		o.oauth2Module.Logout(requestCtx, logoutReqParam, nil))
}

// RevokeClient revokes access of client for the logged-in user
// @Summary      revokes client access
// @Description  It is used by the user in case he/she wants to revoke access for a certain client.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param revokeBody body request_models.RevokeClientBody true "revokeBody"
// @Success      200  {boolean} true
// @Failure      400  {object}  model.ErrorResponse
// @Router       /oauth/revokeClient [post]
func (o *oauth2) RevokeClient(ctx *gin.Context) {
	var revokeRequest request_models.RevokeClientBody
	err := ctx.ShouldBind(&revokeRequest)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	err = o.oauth2Module.RevokeClient(ctx.Request.Context(), revokeRequest)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetAuthorizedClients returns all authorized clients
// @Summary      returns all authorized clients
// @Description  It returns all clients that have resource access other than openid for the logged in user
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Success      200  {object} []dto.AuthorizedClientsResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /oauth/authorizedClients [get]
func (o *oauth2) GetAuthorizedClients(ctx *gin.Context) {
	authorizedClients, err := o.oauth2Module.GetAuthorizedClients(ctx.Request.Context())
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, authorizedClients, nil)
}

// GetOpenIDAuthorizedClients returns all only-openid authorized clients
// @Summary      returns all only-openid authorized clients
// @Description  It returns all clients that have only openid access for the logged-in user
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Success      200  {object} []dto.AuthorizedClientsResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /oauth/openIDAuthorizedClients [get]
func (o *oauth2) GetOpenIDAuthorizedClients(ctx *gin.Context) {
	authorizedClients, err := o.oauth2Module.GetOpenIDAuthorizedClients(ctx.Request.Context())
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, authorizedClients, nil)
}

// UserInfo returns claims about the authenticated End-User
// @Summary      returns claims about the authenticated End-User
// @Description  It returns profile information of user that got logged in using OpenID Connect
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Success      200  {object} []dto.UserInfo
// @Failure      401  {object}  model.ErrorResponse
// @Router       /oauth/userinfo [get]
func (o *oauth2) UserInfo(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	userInfoRsp, err := o.oauth2Module.UserInfo(requestCtx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, userInfoRsp, nil)
}
