package client

import (
	"context"
	"database/sql"

	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
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
		FirstParty:   client.FirstParty,
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

func (c *clientPersistence) GetAllClients(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.Client, *model.MetaData, error) {
	clients, total, err := c.db.GetAllClients(ctx, filters)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no clients found")
			c.logger.Info(ctx, "no clients were found", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading clients")
			c.logger.Error(ctx, "error reading clients", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		}
	}
	clientsDTO := make([]dto.Client, len(clients))
	for k, v := range clients {
		clientsDTO[k] = dto.Client{
			ID:           v.ID,
			Name:         v.Name,
			Status:       v.Status,
			Scopes:       v.Scopes,
			RedirectURIs: utils.StringToArray(v.RedirectUris),
			ClientType:   v.ClientType,
			LogoURL:      v.LogoUrl,
			CreatedAt:    v.CreatedAt,
		}
	}
	return clientsDTO, &model.MetaData{
		FilterParams: filters,
		Total:        total,
		Extra:        nil,
	}, nil
}

func (c *clientPersistence) UpdateClientStatus(ctx context.Context, updateClientStatusParam dto.UpdateClientStatus, clientID uuid.UUID) error {
	_, err := c.db.UpdateClient(ctx, db.UpdateClientParams{
		Status: sql.NullString{String: updateClientStatusParam.Status, Valid: true},
		ID:     clientID,
	})

	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "client not found")
			c.logger.Error(ctx, "error updating client's status", zap.Error(err), zap.Any("client-param", updateClientStatusParam))
			return err
		} else {
			err = errors.ErrUpdateError.Wrap(err, "error updating client status")
			c.logger.Error(ctx, "error updating client's status", zap.Error(err), zap.Any("client-param", updateClientStatusParam))
			return err
		}
	}

	return nil
}

func (c *clientPersistence) UpdateClient(ctx context.Context, client dto.Client) error {
	_, err := c.db.UpdateEntireClient(ctx, db.UpdateEntireClientParams{
		Name:         client.Name,
		LogoUrl:      client.LogoURL,
		ClientType:   client.ClientType,
		RedirectUris: utils.ArrayToString(client.RedirectURIs),
		Scopes:       client.Scopes,
		ID:           client.ID,
	})

	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "client not found")
			c.logger.Error(ctx, "error updating client, ", zap.Error(err), zap.Any("client-param", client))
			return err
		} else {
			err = errors.ErrUpdateError.Wrap(err, "error updating client")
			c.logger.Error(ctx, "error updating client", zap.Error(err), zap.Any("client-param", client))
			return err
		}
	}

	return nil
}
