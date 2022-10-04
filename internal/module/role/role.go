package role

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
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

func (r *roleModule) GetRoleStatus(ctx context.Context, roleName string) (string, error) {
	return r.rolePersistence.GetRoleStatus(ctx, roleName)
}

func (r *roleModule) GetRoleForUser(ctx context.Context, userID string) (string, error) {
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
		r.logger.Warn(ctx, "invalid user id while getting role for user", zap.Error(err), zap.String("user-id", userID))
		return "", err
	}
	return r.rolePersistence.GetRoleForUser(ctx, userIDParsed)
}

func (r *roleModule) GetRoleStatusForUser(ctx context.Context, userID string) (string, error) {
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
		r.logger.Warn(ctx, "invalid user id while getting role for user", zap.Error(err), zap.String("user-id", userID))
		return "", err
	}

	return r.rolePersistence.GetRoleStatusForUser(ctx, userIDParsed)
}
