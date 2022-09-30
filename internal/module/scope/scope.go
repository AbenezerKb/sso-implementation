package scope

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

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
