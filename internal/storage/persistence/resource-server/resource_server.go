package resource_server

import (
	"context"

	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
	"go.uber.org/zap"
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
	resourceServer, err := r.db.CreateResourceServerWithTX(ctx, server)
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create resource server")
		r.logger.Error(ctx, "unable to create resource server", zap.Error(err), zap.Any("server", server))
		return dto.ResourceServer{}, err
	}

	return resourceServer, nil
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

func (r *resourceServerPersistence) GetAllResourceServers(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.ResourceServer, *model.MetaData, error) {
	resourceServers, total, err := r.db.GetAllResourceServers(ctx, filters)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no resource servers found")
			r.logger.Info(ctx, "no resource servers were found", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading resource servers")
			r.logger.Error(ctx, "error reading resource servers", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		}
	}
	return resourceServers, &model.MetaData{
		FilterParams: filters,
		Total:        total,
		Extra:        nil,
	}, nil
}

func (r *resourceServerPersistence) GetResourceServerByID(ctx context.Context, rsID uuid.UUID) (*dto.ResourceServer, error) {
	rs, err := r.db.GetResourceServerByID(ctx, rsID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "resource server not found")
			r.logger.Info(ctx, "resource server was not found", zap.Error(err), zap.String("rs-id", rsID.String()))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read resource server")
			r.logger.Error(ctx, "unable to get resource server by id", zap.Error(err), zap.String("rs-id", rsID.String()))
			return nil, err
		}
	}

	return &dto.ResourceServer{
		ID:        rs.ID,
		Name:      rs.Name,
		CreatedAt: rs.CreatedAt,
		UpdatedAt: rs.UpdatedAt,
		Secret:    rs.Secret,
	}, nil
}
