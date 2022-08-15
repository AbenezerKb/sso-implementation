package client

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"
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
		c.logger.Error(ctx, "couldn't create client", zap.Error(err))
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
