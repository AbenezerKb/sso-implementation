package identity_provider

import (
	"context"

	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
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

	var err error
	ip.ClientSecret, err = utils.Encrypt(ip.ClientSecret, constant.ClientSecretKey)
	if err != nil {
		i.logger.Error(ctx, "error encrypting client secret", zap.Any("client-secret-id", ip.ClientSecret), zap.Error(err))
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

	idPParam.ClientSecret, err = utils.Encrypt(idPParam.ClientSecret, constant.ClientSecretKey)
	if err != nil {
		i.logger.Error(ctx, "error encrypting client secret", zap.Any("idP-id", idPID), zap.Any("client-secret-id", idPParam.ClientSecret), zap.Error(err))
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

func (i *identityProviderModule) GetIdentityProvider(ctx context.Context, idPID string) (dto.IdentityProvider, error) {
	parsedIdPID, err := uuid.Parse(idPID)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid identity provider id")
		i.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("idP-id", idPID))
		return dto.IdentityProvider{}, err
	}

	return i.ipPersistence.GetIdentityProvider(ctx, parsedIdPID)
}

func (i *identityProviderModule) DeleteIdentityProvider(ctx context.Context, idPID string) error {
	parsedIdPID, err := uuid.Parse(idPID)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid identity provider id")
		i.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("idP-id", idPID))
		return err
	}

	err = i.ipPersistence.DeleteIdentityProvider(ctx, parsedIdPID)
	if err != nil {
		return err
	}

	i.logger.Info(ctx, "identity provider deleted", zap.Any("identity-provider-id", parsedIdPID))

	return nil
}

func (i *identityProviderModule) GetAllIdentityProviders(ctx context.Context, filtersQuery db_pgnflt.PgnFltQueryParams) ([]dto.IdentityProvider, *model.MetaData, error) {
	filters, err := filtersQuery.ToFilterParams([]db_pgnflt.FieldType{
		{Name: "name", Type: db_pgnflt.String},
		{Name: "logo_url", Type: db_pgnflt.String},
		{Name: "client_id", Type: db_pgnflt.String},
		{Name: "client_secret", Type: db_pgnflt.String},
		{Name: "redirect_uri", Type: db_pgnflt.String},
		{Name: "authorization_uri", Type: db_pgnflt.String},
		{Name: "token_endpoint_url", Type: db_pgnflt.String},
		{Name: "user_info_endpoint_url", Type: db_pgnflt.String},
		{Name: "status", Type: db_pgnflt.Enum,
			Values: []string{"ACTIVE", "INACTIVE", "PENDING"},
		},
		{Name: "created_at", Type: db_pgnflt.Time},
		{Name: "updated_at", Type: db_pgnflt.Time},
	}, db_pgnflt.Defaults{
		Sort: []db_pgnflt.Sort{
			{
				Field: "created_at",
				Sort:  db_pgnflt.SortDesc,
			},
		},
		PerPage: 10,
	})
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid filter params")
		i.logger.Info(ctx, "invalid filter params were given", zap.Error(err), zap.Any("filters-query", filtersQuery))
		return nil, nil, err
	}
	return i.ipPersistence.GetAllIdentityProviders(ctx, filters)
}
