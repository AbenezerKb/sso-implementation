package role

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/dto"
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
		err := errors.ErrReadError.Wrap(err, "error reading permissions")
		r.logger.Error(ctx, "unable to read permissions", zap.Error(err), zap.String("category", category))
		return nil, err
	}
	if perms == nil {
		err := errors.ErrInvalidUserInput.New(fmt.Sprintf("category %s doesn't exist", category))
		r.logger.Info(ctx, "category was not found", zap.String("category", category), zap.Error(err))
		return nil, err
	}
	return perms, nil
}

func (r *rolePersistence) GetRoleStatus(ctx context.Context, roleName string) (string, error) {
	status, err := r.db.GetRoleStatus(ctx, roleName)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return "", nil
		}
		err := errors.ErrReadError.Wrap(err, "error fetching role status")
		r.logger.Error(ctx, "unable to fetch role status", zap.Error(err), zap.String("role-name", roleName))
		return "", err
	}
	return status.String, nil
}

func (r *rolePersistence) GetRoleForUser(ctx context.Context, userID uuid.UUID) (string, error) {
	role, err := r.db.GetRoleForUser(ctx, userID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return "", nil
		}
		err := errors.ErrReadError.Wrap(err, "error fetching role for user")
		r.logger.Error(ctx, "error while reading role for user", zap.Error(err), zap.Any("user-id", userID))
		return "", err
	}

	return role, nil
}

func (r *rolePersistence) GetRoleStatusForUser(ctx context.Context, userID uuid.UUID) (string, error) {
	role, err := r.GetRoleForUser(ctx, userID)
	if err != nil {
		return "", err
	}

	status, err := r.GetRoleStatus(ctx, role)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *rolePersistence) CreateRole(ctx context.Context, role dto.Role) (dto.Role, error) {
	roleSaved, err := r.db.CreateRoleTX(ctx, role.Name, role.Permissions)
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "error creating role")
		r.logger.Error(ctx, "error while creating a role", zap.Error(err), zap.Any("role", role))
		return dto.Role{}, err
	}

	return roleSaved, nil
}

func (r *rolePersistence) CheckIfPermissionExists(ctx context.Context, perm string) (bool, error) {
	exist, err := r.db.CheckIfPermissionExists(ctx, perm)
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "error checking if permission exists")
		r.logger.Error(ctx, "error while checking for permission existence", zap.Error(err), zap.Any("permission", perm))
	}

	return exist, nil
}