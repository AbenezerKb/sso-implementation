package role

import (
	"context"
	"sso/internal/constant/permissions"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
)

type roleModule struct {
	logger          logger.Logger
	rolePersistence storage.RolePersistence
}

func InitRole(logger logger.Logger, rolePersistence storage.RolePersistence) module.RoleModule {
	return &roleModule{
		logger:          logger,
		rolePersistence: rolePersistence,
	}
}

func (r *roleModule) GetAllPermissions(ctx context.Context, category string) ([]permissions.Permission, error) {
	return r.rolePersistence.GetAllPermissions(ctx, category)
}
