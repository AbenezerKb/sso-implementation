package role

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
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
		perms, err := r.db.GetAllPermissions(ctx)
		if err != nil {
			err := errors.ErrReadError.Wrap(err, "error reading permissions")
			r.logger.Error(ctx, "unable to read permissions", zap.Error(err))
			return nil, err
		}
		return perms, nil
	}
	perms, err := r.db.GetPermissionsOfCategory(ctx, category)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrInvalidUserInput.Wrap(err, fmt.Sprintf("category %s doesn't exist", category))
			r.logger.Info(ctx, "category was not found", zap.String("category", category), zap.Error(err))
			return nil, err
		}
		err := errors.ErrReadError.Wrap(err, "error reading permissions", zap.String("category", category))
		r.logger.Error(ctx, "unable to read permissions", zap.Error(err), zap.String("category", category))
		return nil, err
	}
	return perms, nil
}
