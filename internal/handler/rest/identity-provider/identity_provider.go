package identity_provider

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
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
// @Param role body dto.IdentityProvider true "identityProvider"
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
