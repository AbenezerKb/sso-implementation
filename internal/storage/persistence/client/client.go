package client

import (
	"context"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
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

func (c *clientPersistence) Create(ctx context.Context, client dto.Client) (*dto.Client, error) {
	return nil, nil
}
