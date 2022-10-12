package identity_provider

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (i *identityProviderModule) UpdateIdentityProvider(ctx context.Context, idPParam dto.IdentityProvider, idPID string) error {
	parsedIdPID, err := uuid.Parse(idPID)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid identity provider id")
		i.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("idP-id", idPID))
		return err
	}
	if err := idPParam.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		i.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("identity-provider", idPParam))
		return err
	}

	idPParam.ID = parsedIdPID

	err = i.ipPersistence.UpdateIdentityProvider(ctx, idPParam)
	if err != nil {
		return err
	}

	i.logger.Info(ctx, "identity provider updated", zap.Any("identity-provider-id", parsedIdPID), zap.Any("updated-to", idPParam))
	return nil
}

func (i *identityProviderModule) GetIdentityProvider(ctx context.Context, idPID string) (*dto.IdentityProvider, error) {
	parsedIdPID, err := uuid.Parse(idPID)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid identity provider id")
		i.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("idP-id", idPID))
		return nil, err
	}

	return i.ipPersistence.GetIdentityProvider(ctx, parsedIdPID)
}