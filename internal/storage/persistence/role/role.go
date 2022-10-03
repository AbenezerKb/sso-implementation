package role

import (
	"context"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/constant/permissions"
	"sso/internal/storage"
	"sso/platform/logger"
)

type rolePersistence struct {
	logger logger.Logger
	db     *persistencedb.PersistenceDB
}

func InitRolePersistence(logger logger.Logger, db *persistencedb.PersistenceDB) storage.RolePersistence {
	return &rolePersistence{
		logger: logger,
		db:     db,
	}
}

func (r *rolePersistence) GetAllPermissions(ctx context.Context, category string) ([]permissions.Permission, error) {
	if category == "" {
		return r.db.GetAllPermissions(ctx)
	}
	return r.db.GetPermissionsOfCategory(ctx, category)
}
