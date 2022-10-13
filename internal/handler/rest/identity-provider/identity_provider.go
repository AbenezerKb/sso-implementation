package identity_provider

import (
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type identityProvider struct {
	logger   logger.Logger
	ipModule module.IdentityProviderModule
}

func InitIdentityProvider(logger logger.Logger, ipModule module.IdentityProviderModule) rest.IdentityProvider {
	return &identityProvider{
		logger:   logger,
		ipModule: ipModule,
	}
}

// CreateIdentityProvider is used to create an identity provider
// @Summary create identity provider
// @Description create an identity provider
// @ID create-identity-provider
// @Tags identityProvider
// @Accept  json
// @Produce  json
// @Param identityProvider body dto.IdentityProvider true "identityProvider"
// @Success 200 {object} dto.IdentityProvider
// @Failure 400 {object} model.ErrorResponse
// @Router /identityProviders [post]
// @Security BearerAuth
func (i *identityProvider) CreateIdentityProvider(ctx *gin.Context) {
	var ip dto.IdentityProvider
	err := ctx.ShouldBind(&ip)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		i.logger.Info(ctx, "could not bind to dto.IdentityProvider. invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	ipCreated, err := i.ipModule.CreateIdentityProvider(requestCtx, ip)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusCreated, ipCreated, nil)
}

// UpdateIdentityProvider is used to update an identity provider
// @Summary update identity provider
// @Description update an identity provider
// @ID update-identity-provider
// @Tags identityProvider
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Param identityProvider body dto.IdentityProvider true "identityProvider"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.ErrorResponse
// @Router /identityProviders/{id} [put]
// @Security BearerAuth
func (i *identityProvider) UpdateIdentityProvider(ctx *gin.Context) {
	idPID := ctx.Param("id")

	var idPParam dto.IdentityProvider
	err := ctx.ShouldBind(&idPParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		i.logger.Info(ctx, "could not bind to dto.IdentityProvider. invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	requestCtx := ctx.Request.Context()
	err = i.ipModule.UpdateIdentityProvider(requestCtx, idPParam, idPID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetIdentityProvider is used to get a particular identity provider
// @Summary get identity provider
// @Description get an identity provider
// @ID get-identity-provider
// @Tags identityProvider
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} dto.IdentityProvider
// @Failure 400 {object} model.ErrorResponse
// @Router /identityProviders/{id} [get]
// @Security BearerAuth
func (i *identityProvider) GetIdentityProvider(ctx *gin.Context) {
	idPID := ctx.Param("id")

	requestCtx := ctx.Request.Context()
	idP, err := i.ipModule.GetIdentityProvider(requestCtx, idPID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, idP, nil)
}

// DeleteIdentityProvider is used to delete a particular identity provider
// @Summary delete identity provider
// @Description delete an identity provider
// @ID delete-identity-provider
// @Tags identityProvider
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 204 {object} model.Response
// @Failure 400 {object} model.ErrorResponse
// @Router /identityProviders/{id} [delete]
// @Security BearerAuth
func (i *identityProvider) DeleteIdentityProvider(ctx *gin.Context) {
	idPID := ctx.Param("id")

	requestCtx := ctx.Request.Context()
	err := i.ipModule.DeleteIdentityProvider(requestCtx, idPID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusNoContent, nil, nil)
}

// GetIdentityProviders is used to get a all identity providers
// @Summary get all identity providers
// @Description get all identity providers
// @ID get-identity-providers
// @Tags identityProvider
// @Accept  json
// @Produce  json
// @Success 200 {object} []dto.IdentityProvider
// @Failure 400 {object} model.ErrorResponse
// @Router /identityProviders [get]
// @Security BearerAuth
func (i *identityProvider) GetAllIdentityProviders(ctx *gin.Context) {
	var filtersParam request_models.PgnFltQueryParams
	err := ctx.BindQuery(&filtersParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		i.logger.Info(ctx, "invalid query params", zap.Error(err), zap.Any("query-params", ctx.Request.URL.Query()))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	idPs, metaData, err := i.ipModule.GetAllIdentityProviders(requestCtx, filtersParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, idPs, metaData)
}
