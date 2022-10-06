package role

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
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

func (r *roleModule) CreateRole(ctx context.Context, role dto.Role) (dto.Role, error) {
	if err := role.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Any("role", role))
		return dto.Role{}, err
	}

	// check for invalid permissions
	for i := 0; i < len(role.Permissions); i++ {
		exists, err := r.rolePersistence.CheckIfPermissionExists(ctx, role.Permissions[i])
		if err != nil {
			return dto.Role{}, err
		}
		if !exists {
			err := errors.ErrInvalidUserInput.New(fmt.Sprintf("permission %s doesn't exist", role.Permissions[i]))
			r.logger.Info(ctx, "permission doesn't exist")
			return dto.Role{}, err
		}
	}

	return r.rolePersistence.CreateRole(ctx, role)
}

func (r *roleModule) GetAllRoles(ctx context.Context, filtersQuery request_models.PgnFltQueryParams) ([]dto.Role, *model.MetaData, error) {
	filters, err := filtersQuery.ToFilterParams(dto.Role{})
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid filter params")
		r.logger.Info(ctx, "invalid filter params were given", zap.Error(err), zap.Any("filters-query", filtersQuery))
		return nil, nil, err
	}
	return r.rolePersistence.GetAllRoles(ctx, filters)
}
