package client

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type clientPersistence struct {
	logger logger.Logger
	db     *db.Queries
}

func InitClient(log logger.Logger, db *db.Queries) storage.ClientPersistence {
	return &clientPersistence{
		logger: log,
		db:     db,
	}
}

func (c *clientPersistence) Create(ctx context.Context, clientParam dto.Client) (*dto.Client, error) {
	client, err := c.db.CreateClient(ctx, db.CreateClientParams{
		Name:         clientParam.Name,
		ClientType:   clientParam.ClientType,
		RedirectUris: utils.ArrayToString(clientParam.RedirectURIs),
		Scopes:       clientParam.Scopes,
		Secret:       clientParam.Secret,
		LogoUrl:      clientParam.LogoURL,
	})
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "couldn't create client")
		c.logger.Error(ctx, "couldn't create client", zap.Error(err), zap.Any("client", clientParam))
		return nil, err
	}
	return &dto.Client{
		ID:           client.ID,
		Name:         client.Name,
		ClientType:   client.ClientType,
		RedirectURIs: utils.StringToArray(client.RedirectUris),
		Scopes:       client.Scopes,
		Secret:       client.Secret,
		LogoURL:      client.LogoUrl,
		Status:       client.Status,
	}, nil
}

func (c *clientPersistence) GetClientByID(ctx context.Context, id uuid.UUID) (*dto.Client, error) {
	client, err := c.db.GetClientByID(ctx, id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no client found")
			c.logger.Info(ctx, "client not found", zap.Error(err), zap.Any("client-id", id))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading the client")
			c.logger.Error(ctx, "error reading the client", zap.Error(err), zap.Any("client-id", id))
			return nil, err
		}
	}

	return &dto.Client{
		ID:           client.ID,
		Name:         client.Name,
		Status:       client.Status,
		Secret:       client.Secret,
		Scopes:       client.Scopes,
		RedirectURIs: utils.StringToArray(client.RedirectUris),
		ClientType:   client.ClientType,
		LogoURL:      client.LogoUrl,
	}, nil

}

func (c *clientPersistence) DeleteClientByID(ctx context.Context, id uuid.UUID) error {
	_, err := c.db.DeleteClient(ctx, id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "client not found")
			c.logger.Info(ctx, "client not found", zap.Error(err), zap.Any("client-id", id))
			return err
		}
		err = errors.ErrDBDelError.Wrap(err, "error deleting client")
		c.logger.Error(ctx, "error deleting client", zap.Error(err), zap.Any("client-id", id))
		return err

	}

	return nil
}
