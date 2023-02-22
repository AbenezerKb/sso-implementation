package client

import (
	"context"

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

type clientModule struct {
	logger            logger.Logger
	clientPersistence storage.ClientPersistence
}

func InitClient(log logger.Logger, clientPersistence storage.ClientPersistence) module.ClientModule {
	return &clientModule{
		logger:            log,
		clientPersistence: clientPersistence,
	}
}

func (c *clientModule) Create(ctx context.Context, clientParam dto.Client) (*dto.Client, error) {
	if err := clientParam.ValidateClient(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	// TODO: check scope on the resource server
	clientParam.Secret = utils.GenerateRandomString(25, true)

	return c.clientPersistence.Create(ctx, clientParam)
}

func (c *clientModule) GetClientByID(ctx context.Context, id string) (*dto.Client, error) {
	clientID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid client id")
		c.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("client-id", id))
		return nil, err
	}
	return c.clientPersistence.GetClientByID(ctx, clientID)
}

func (c *clientModule) GetAllClients(ctx context.Context, filtersQuery db_pgnflt.PgnFltQueryParams) ([]dto.Client, *model.MetaData, error) {
	filters, err := filtersQuery.ToFilterParams([]db_pgnflt.FieldType{
		{Name: "name", Type: db_pgnflt.String},
		{Name: "client_type", Type: db_pgnflt.Enum, Values: []string{"Confidential", "Public"}},
		{Name: "redirect_uris", Type: db_pgnflt.String},
		{Name: "scopes", Type: db_pgnflt.String},
		{Name: "status", Type: db_pgnflt.Enum,
			Values: []string{"ACTIVE", "INACTIVE", "PENDING"},
		},
		{Name: "created_at", Type: db_pgnflt.Time},
		{Name: "first_party", Type: db_pgnflt.Boolean},
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
		c.logger.Info(ctx, "invalid filter params were given", zap.Error(err), zap.Any("filters-query", filtersQuery))
		return nil, nil, err
	}
	return c.clientPersistence.GetAllClients(ctx, filters)
}
func (c *clientModule) DeleteClientByID(ctx context.Context, id string) error {
	clientID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "client not found")
		c.logger.Info(ctx, "parse error", zap.Error(err), zap.String("client-id", id))
		return err
	}

	// TODO: before deleting client we should de something about rf token issued to this client
	// TODO: before deleting this client we should do something about the auth_histories of this client
	err = c.clientPersistence.DeleteClientByID(ctx, clientID)
	if err != nil {
		return err
	}

	return nil
}

func (c *clientModule) UpdateClientStatus(ctx context.Context, updateClientStatusParam dto.UpdateClientStatus, id string) error {
	clientID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid client id")
		c.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return err
	}

	if err := updateClientStatusParam.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.logger.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	err = c.clientPersistence.UpdateClientStatus(ctx, updateClientStatusParam, clientID)
	if err != nil {
		return err
	}
	return nil
}

func (c *clientModule) UpdateClient(ctx context.Context, client dto.Client, id string) error {
	clientID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "invalid client id")
		c.logger.Error(ctx, "parse error", zap.Error(err), zap.Any("client-id", id))
		return err
	}

	if err := client.ValidateClient(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.logger.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	client.ID = clientID

	return c.clientPersistence.UpdateClient(ctx, client)
}
