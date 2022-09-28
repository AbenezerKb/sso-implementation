package resource_server

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
)

type resourceServerPersistence struct {
	logger logger.Logger
	db     *db.Queries
}

func InitResourceServerPersistence(logger logger.Logger, db *db.Queries) storage.ResourceServerPersistence {
	return &resourceServerPersistence{
		logger: logger,
		db:     db,
	}
}

func (r *resourceServerPersistence) CreateResourceServer(ctx context.Context, server dto.ResourceServer) (dto.ResourceServer, error) {
	resourceServer, err := r.db.CreateResourceServer(ctx, server.Name)
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create resource server")
		r.logger.Error(ctx, "unable to create resource server", zap.Error(err), zap.Any("server", server))
		return dto.ResourceServer{}, err
	}
	return dto.ResourceServer{
		ID:        resourceServer.ID,
		Name:      resourceServer.Name,
		CreatedAt: resourceServer.CreatedAt,
		UpdatedAt: resourceServer.UpdatedAt,
		Scopes:    nil,
	}, nil
}
