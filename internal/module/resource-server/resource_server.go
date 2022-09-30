package resource_server

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
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
