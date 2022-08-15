package client

import (
	"context"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
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

func (c *clientModule) Create(ctx context.Context, client dto.Client) (*dto.Client, error) {
	return nil, nil
}
