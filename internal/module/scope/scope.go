package scope

import (
	"context"

	"sso/internal/constant/errors"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
	"go.uber.org/zap"
)

type scopeModule struct {
	logger           logger.Logger
	scopePersistence storage.ScopePersistence
}

func InitScope(logger logger.Logger, scopePersistence storage.ScopePersistence) module.ScopeModule {
	return &scopeModule{
		logger:           logger,
		scopePersistence: scopePersistence,
	}
}

func (s *scopeModule) GetScope(ctx context.Context, scope string) (dto.Scope, error) {
	return s.scopePersistence.GetScope(ctx, scope)
}

func (s *scopeModule) CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error) {
	// validate the scope
	if err := scope.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.logger.Info(ctx, "invalid input", zap.Error(err))
		return dto.Scope{}, err
	}

	name := scope.Name
	if scope.ResourceServerName != "" {
		// TODO: check if resource server with this name exists
		// TODO: check if scope name is unique for that resource server
		name = scope.ResourceServerName + "." + name
	}
	// TODO: check if scope name is unique for non-resource-server scopes

	return s.scopePersistence.CreateScope(ctx, dto.Scope{
		Name:               name,
		Description:        scope.Description,
		ResourceServerName: scope.ResourceServerName,
	})
}

func (s *scopeModule) GetAllScopes(ctx context.Context, filtersQuery db_pgnflt.PgnFltQueryParams) ([]dto.Scope, *model.MetaData, error) {
	filters, err := filtersQuery.ToFilterParams([]db_pgnflt.FieldType{
		{Name: "name", Type: db_pgnflt.String},
		{Name: "description", Type: db_pgnflt.String},
		{Name: "resource_server_name", Type: db_pgnflt.String},
		{Name: "status", Type: db_pgnflt.Enum,
			Values: []string{"ACTIVE", "INACTIVE", "PENDING"},
		},
		{Name: "created_at", Type: db_pgnflt.Time},
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
		s.logger.Info(ctx, "invalid filter params were given", zap.Error(err), zap.Any("filters-query", filtersQuery))
		return nil, nil, err
	}
	return s.scopePersistence.GetAllScopes(ctx, filters)
}

func (s *scopeModule) DeleteScopeByName(ctx context.Context, name string) error {
	return s.scopePersistence.DeleteScopeByName(ctx, name)
}

func (s *scopeModule) UpdateScope(ctx context.Context, updateScopeParam dto.UpdateScopeParam, scopeName string) error {
	if err := updateScopeParam.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.logger.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	return s.scopePersistence.UpdateScope(ctx, dto.Scope{
		Name:        scopeName,
		Description: updateScopeParam.Description,
	})
}
