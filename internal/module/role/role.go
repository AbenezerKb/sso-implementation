package role

import (
	"context"
	"fmt"
	"sso/internal/constant/errors"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (r *roleModule) GetAllPermissions(ctx context.Context, category string) ([]dto.Permission, error) {
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

func (r *roleModule) UpdateRoleStatus(ctx context.Context, updateRoleStatusParam dto.UpdateRoleStatus, roleName string) error {
	if err := updateRoleStatusParam.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	err := r.rolePersistence.UpdateRoleStatus(ctx, updateRoleStatusParam, roleName)
	if err != nil {
		return err
	}
	return nil
}

func (r *roleModule) GetRoleByName(ctx context.Context, roleName string) (dto.Role, error) {
	return r.rolePersistence.GetRoleByName(ctx, roleName)
}

func (r *roleModule) DeleteRole(ctx context.Context, roleName string) error {
	return r.rolePersistence.DeleteRole(ctx, roleName)
}

func (r *roleModule) UpdateRole(ctx context.Context, updateRole dto.UpdateRole) (dto.Role, error) {
	if err := updateRole.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input while updating role", zap.Error(err), zap.Any("role", updateRole))
		return dto.Role{}, err
	}

	// check for invalid permissions
	for i := 0; i < len(updateRole.Permissions); i++ {
		exists, err := r.rolePersistence.CheckIfPermissionExists(ctx, updateRole.Permissions[i])
		if err != nil {
			return dto.Role{}, err
		}
		if !exists {
			err := errors.ErrInvalidUserInput.New(fmt.Sprintf("permission %s does not exist", updateRole.Permissions[i]))
			r.logger.Info(ctx, "permission doesn't exist")
			return dto.Role{}, err
		}
	}
	return r.rolePersistence.UpdateRole(ctx, updateRole)
}
