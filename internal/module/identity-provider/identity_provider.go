package identity_provider

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
)

type identityProviderModule struct {
	logger        logger.Logger
	ipPersistence storage.IdentityProviderPersistence
}

func InitIdentityProvider(logger logger.Logger, ipPersistence storage.IdentityProviderPersistence) module.IdentityProviderModule {
	return &identityProviderModule{
		logger:        logger,
		ipPersistence: ipPersistence,
	}
}

func (i *identityProviderModule) CreateIdentityProvider(ctx context.Context, ip dto.IdentityProvider) (dto.IdentityProvider, error) {
	if err := ip.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		i.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("identity-provider", ip))
		return dto.IdentityProvider{}, err
	}

	return i.ipPersistence.CreateIdentityProvider(ctx, ip)
}
