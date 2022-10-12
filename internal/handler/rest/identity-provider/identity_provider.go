package identity_provider

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

	ipCreated, err := i.ipModule.CreateIdentityProvider(ctx, ip)
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
// @Param identityProvider body dto.IdentityProvider true "identityProvider"
// @Success 200 {object} dto.IdentityProvider
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

	err = i.ipModule.UpdateIdentityProvider(ctx, idPParam, idPID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
