package resource_server

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"
)

type resourceServerPersistence struct {
	logger logger.Logger
	db     *persistencedb.PersistenceDB
}

func InitResourceServerPersistence(logger logger.Logger, db *persistencedb.PersistenceDB) storage.ResourceServerPersistence {
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

func (r *resourceServerPersistence) GetResourceServerByName(ctx context.Context, name string) (dto.ResourceServer, error) {
	resourceServer, err := r.db.GetResourceServerByName(ctx, name)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "resource server not found")
			r.logger.Info(ctx, "resource server was not found", zap.Error(err), zap.String("resource-server-name", name))
			return dto.ResourceServer{}, err
		}
		err = errors.ErrReadError.Wrap(err, "could not read the resource server")
		r.logger.Error(ctx, "unable to read the resource server", zap.Error(err), zap.Any("resource-server-name", name))
		return dto.ResourceServer{}, err
	}

	return dto.ResourceServer{
		ID:        resourceServer.ID,
		Name:      resourceServer.Name,
		CreatedAt: resourceServer.CreatedAt,
		UpdatedAt: resourceServer.UpdatedAt,
	}, nil
}
