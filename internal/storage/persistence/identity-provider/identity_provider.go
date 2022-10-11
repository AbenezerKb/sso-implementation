package identity_provider

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"
)

type identityProviderPersistence struct {
	logger logger.Logger
	db     *persistencedb.PersistenceDB
}

func InitIdentityProviderPersistence(logger logger.Logger, db *persistencedb.PersistenceDB) storage.IdentityProviderPersistence {
	return &identityProviderPersistence{
		logger: logger,
		db:     db,
	}
}

func (i *identityProviderPersistence) CreateIdentityProvider(ctx context.Context, ip dto.IdentityProvider) (dto.IdentityProvider, error) {
	ipDB, err := i.db.CreateIdentityProvider(ctx, db.CreateIdentityProviderParams{
		Name: ip.Name,
		LogoUrl: sql.NullString{
			String: ip.LogoURI,
			Valid:  true,
		},
		ClientID:         ip.ClientID,
		ClientSecret:     ip.ClientSecret,
		RedirectUri:      ip.RedirectURI,
		AuthorizationUri: ip.AuthorizationURI,
		TokenEndpointUrl: ip.TokenEndpointURI,
		UserInfoEndpointUrl: sql.NullString{
			String: ip.UserInfoEndpointURI,
			Valid:  true,
		},
	})

	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "error creating identity provider")
		i.logger.Error(ctx, "error while creating identity provider", zap.Error(err), zap.Any("identity-provider", ip))
		return dto.IdentityProvider{}, err
	}

	return dto.IdentityProvider{
		ID:                  ipDB.ID,
		Name:                ipDB.Name,
		LogoURI:             ipDB.LogoUrl.String,
		ClientID:            ipDB.ClientID,
		ClientSecret:        ipDB.ClientSecret,
		RedirectURI:         ipDB.RedirectUri,
		AuthorizationURI:    ipDB.AuthorizationUri,
		TokenEndpointURI:    ipDB.TokenEndpointUrl,
		UserInfoEndpointURI: ipDB.UserInfoEndpointUrl.String,
		Status:              ipDB.Status.String,
		CreatedAt:           ipDB.CreatedAt,
		UpdatedAt:           ipDB.UpdatedAt,
	}, nil
}
