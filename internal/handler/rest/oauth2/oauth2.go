package oauth2

import (
	"net/http"
	"net/url"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/state"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
	"strings"

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
	errorURL, err := url.Parse(state.ConsentURL)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	errQuery := errorURL.Query()

	authRequestParam := dto.AuthorizationRequestParam{}
	err = ctx.ShouldBindQuery(&authRequestParam)
	if err != nil {
		o.logger.Info(ctx, zap.Error(err).String)
		errQuery.Set("error", "invalid_request")
		errQuery.Set("error_description", err.Error())
		errQuery.Set("error_code", "400")
		errorURL.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusBadRequest, errorURL.String())
		return
	}

	// errRedirectURI, err := url.Parse(authRequestParam.RedirectURI)
	// if err != nil {
	// 	_ = ctx.Error(err)
	// 	return
	// }
	authRequestParam.ClientID, err = uuid.Parse(ctx.Query("client_id"))
	if err != nil {
		o.logger.Info(ctx, zap.Error(err).String)
		errQuery.Set("error", "invalid_client_id")
		errQuery.Set("error_description", "invalid client id.")
		errQuery.Set("code", "400")

		errorURL.RawQuery = errQuery.Encode()
		ctx.Redirect(http.StatusFound, errorURL.String())
		return
	}

	consentId, authErrRsp, err := o.oauth2Module.Authorize(ctx.Request.Context(), authRequestParam)
	if err != nil {
		errQuery.Set("error", authErrRsp.Error)
		errQuery.Set("error_description", authErrRsp.ErrorDescription)
		errQuery.Set("error_code", "400")

		errorURL.RawQuery = errQuery.Encode()

		ctx.Redirect(http.StatusFound, errorURL.String())
		return
	}

	consentURL, err := url.Parse(state.ConsentURL)
	if err != nil {
		_ = ctx.Error(err)
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
// @Success      200  {object}  dto.ConsentData
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Router       /oauth/consent/{id} [get]
func (o *oauth2) GetConsentByID(ctx *gin.Context) {
	consentID := ctx.Param("id")
	userID := ctx.Query("user_id")

	consent, err := o.oauth2Module.GetConsentByID(ctx.Request.Context(), consentID, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, consent)
}

// Approval is used to approve consent.
// @Summary      Approval.
// @Description  is used to approve consent.
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @param consentId query string true "consentId"
// @param access query string true "access"
// @success 	 200
// @Failure      400  {object}  model.ErrorResponse "invalid input"
// @Header       200,400            {string}  Location  "redirect_uri"
// @Router       /oauth/approval [get]
func (o *oauth2) Approval(ctx *gin.Context) {
	consentId := ctx.Query("consentId")
	accessRqResult := ctx.Query("access")
	failureReason := ctx.Query("failureReason")
	// userID := ctx.GetString("user_id")
	if consentId == "" || accessRqResult == "" {
		o.logger.Error(ctx, "invalid input", zap.String("phone", consentId), zap.String("access", accessRqResult))
		_ = ctx.Error(errors.ErrInvalidUserInput.New("invalid input"))
		return
	}
	consent, err := o.oauth2Module.Approval(ctx.Request.Context(), consentId, accessRqResult)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	redirectURI, err := url.Parse(consent.RedirectURI)
	if err != nil {
		o.logger.Error(ctx, "invalid input", zap.String("redirect_uri", consent.RedirectURI))
		_ = ctx.Error(err)
		return
	}
	query := redirectURI.Query()
	query.Set("state", consent.State)

	if accessRqResult == "true" {
		authCode, st, err := o.oauth2Module.IssueAuthCode(ctx, consent)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		query.Set("code", authCode)
		if st != "" {
			query.Set("state", st)
		}
		if strings.Contains(consent.ResponseType, "id_token") {
			query.Set("id_token", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")
		}

	} else {
		if failureReason == "" {
			failureReason = "access_denied"
		}
		query.Set("error", failureReason)
	}

	redirectURI.RawQuery = query.Encode()
	ctx.Redirect(http.StatusFound, redirectURI.String())
}
