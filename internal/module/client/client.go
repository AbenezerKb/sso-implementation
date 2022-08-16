package client

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"
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
