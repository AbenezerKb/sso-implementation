package identity_provider

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (i *identityProviderPersistence) UpdateIdentityProvider(ctx context.Context, idPParam dto.IdentityProvider) error {
	_, err := i.db.UpdateIdentityProvider(ctx, db.UpdateIdentityProviderParams{
		Name:                idPParam.Name,
		LogoUrl:             sql.NullString{String: idPParam.LogoURI, Valid: true},
		ClientID:            idPParam.ClientID,
		ClientSecret:        idPParam.ClientSecret,
		RedirectUri:         idPParam.RedirectURI,
		AuthorizationUri:    idPParam.AuthorizationURI,
		TokenEndpointUrl:    idPParam.TokenEndpointURI,
		UserInfoEndpointUrl: sql.NullString{String: idPParam.UserInfoEndpointURI, Valid: true},
		ID:                  idPParam.ID,
	})

	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "identity provider not found")
			i.logger.Error(ctx, "error updating identity provider, ", zap.Error(err), zap.Any("idP-param", idPParam))
			return err
		} else {
			err = errors.ErrUpdateError.Wrap(err, "error updating identity provider")
			i.logger.Error(ctx, "error updating identity provider", zap.Error(err), zap.Any("idP-param", idPParam))
			return err
		}
	}

	return nil
}

func (i *identityProviderPersistence) GetIdentityProvider(ctx context.Context, idPID uuid.UUID) (*dto.IdentityProvider, error) {
	idP, err := i.db.GetIdentityProvider(ctx, idPID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no identity provider found")
			i.logger.Info(ctx, "identity provider not found", zap.Error(err), zap.Any("idP-id", idPID))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading the identity provider")
			i.logger.Error(ctx, "error reading the identity provider", zap.Error(err), zap.Any("idP-id", idPID))
			return nil, err
		}
	}

	return &dto.IdentityProvider{
		ID:                  idP.ID,
		Name:                idP.Name,
		LogoURI:             idP.LogoUrl.String,
		ClientID:            idP.ClientID,
		ClientSecret:        idP.ClientSecret,
		RedirectURI:         idP.RedirectUri,
		AuthorizationURI:    idP.AuthorizationUri,
		TokenEndpointURI:    idP.TokenEndpointUrl,
		UserInfoEndpointURI: idP.UserInfoEndpointUrl.String,
		Status:              idP.Status.String,
		CreatedAt:           idP.CreatedAt,
	}, nil
}

func (i *identityProviderPersistence) DeleteIdentityProvider(ctx context.Context, idPID uuid.UUID) error {
	_, err := i.db.DeleteIdentityProvider(ctx, idPID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "identity provider not found")
			i.logger.Info(ctx, "identity provider not found", zap.Error(err), zap.Any("idP-id", idPID))
			return err
		} else {
			err = errors.ErrDBDelError.Wrap(err, "error deleting the identity provider")
			i.logger.Error(ctx, "error deleting the identity provider", zap.Error(err), zap.Any("idP-id", idPID))
			return err
		}
	}

	return nil
}

func (i *identityProviderPersistence) GetAllIdentityProviders(ctx context.Context, filters request_models.FilterParams) ([]dto.IdentityProvider, *model.MetaData, error) {
	idPs, total, err := i.db.GetAllIdentityProviders(ctx, utils.ComposeFilterSQL(ctx, filters, i.logger))
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no identity providers found")
			i.logger.Info(ctx, "no identity provider were found", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading identity providers")
			i.logger.Error(ctx, "error reading identity providers", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		}
	}
	idpPsDTO := make([]dto.IdentityProvider, len(idPs))
	for k, v := range idPs {
		idpPsDTO[k] = dto.IdentityProvider{
			ID:                  v.ID,
			Name:                v.Name,
			Status:              v.Status.String,
			LogoURI:             v.LogoUrl.String,
			ClientSecret:        v.ClientSecret,
			ClientID:            v.ClientID,
			RedirectURI:         v.RedirectUri,
			AuthorizationURI:    v.AuthorizationUri,
			TokenEndpointURI:    v.TokenEndpointUrl,
			UserInfoEndpointURI: v.UserInfoEndpointUrl.String,
			CreatedAt:           v.CreatedAt,
		}
	}
	return idpPsDTO, &model.MetaData{
		FilterParams: filters,
		Total:        total,
		Extra:        nil,
	}, nil
}
