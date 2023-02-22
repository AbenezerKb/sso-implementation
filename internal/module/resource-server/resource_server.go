package resource_server

import (
	"context"

	"sso/internal/constant/errors"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
	"go.uber.org/zap"
)

type resourceServerModule struct {
	logger                    logger.Logger
	resourceServerPersistence storage.ResourceServerPersistence
	scopePersistence          storage.ScopePersistence
}

func InitResourceServer(logger logger.Logger, rsp storage.ResourceServerPersistence, sp storage.ScopePersistence) module.ResourceServerModule {
	return &resourceServerModule{
		logger:                    logger,
		resourceServerPersistence: rsp,
		scopePersistence:          sp,
	}
}

func (r *resourceServerModule) CreateResourceServer(ctx context.Context, server dto.ResourceServer) (dto.ResourceServer, error) {
	// validate fields
	if err := server.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err))
		return dto.ResourceServer{}, err
	}

	// check if resource server name is unique
	_, err := r.resourceServerPersistence.GetResourceServerByName(ctx, server.Name)
	if err == nil {
		err := errors.ErrInvalidUserInput.New("this server name is taken")
		r.logger.Info(ctx, "invalid input", zap.Error(err))
		return dto.ResourceServer{}, err
	}

	// append server name to scope names
	for i := 0; i < len(server.Scopes); i++ {
		server.Scopes[i].Name = server.Name + "." + server.Scopes[i].Name
	}

	// create resource server
	return r.resourceServerPersistence.CreateResourceServer(ctx, server)
}

func (r *resourceServerModule) GetAllResourceServers(ctx context.Context, filtersQuery db_pgnflt.PgnFltQueryParams) ([]dto.ResourceServer, *model.MetaData, error) {
	filters, err := filtersQuery.ToFilterParams([]db_pgnflt.FieldType{
		{Name: "name", Type: db_pgnflt.String, DBName: "rs.name"},
		{Name: "created_at", Type: db_pgnflt.Time, DBName: "rs.created_at"},
		{Name: "updated_at", Type: db_pgnflt.Time, DBName: "rs.updated_at"},
		{Name: "name", Type: db_pgnflt.String, DBName: "sc.name"},
		{Name: "description", Type: db_pgnflt.String, DBName: "sc.description"},
		{Name: "status", Type: db_pgnflt.Enum,
			Values: []string{"ACTIVE", "INACTIVE", "PENDING"},
			DBName: "sc.status",
		},
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
		r.logger.Info(ctx, "invalid filter params were given", zap.Error(err), zap.Any("filters-query", filtersQuery))
		return nil, nil, err
	}
	return r.resourceServerPersistence.GetAllResourceServers(ctx, filters)
}

func (r *resourceServerModule) GetResourceServerByID(ctx context.Context, rsID string) (*dto.ResourceServer, error) {
	userID, err := uuid.Parse(rsID)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "resource server not found")
		r.logger.Info(ctx, "parse error", zap.Error(err), zap.String("resource-server-id", rsID))
		return nil, err
	}

	return r.resourceServerPersistence.GetResourceServerByID(ctx, userID)
}
