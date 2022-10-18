package user

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userPersistence struct {
	logger logger.Logger
	db     *persistencedb.PersistenceDB
}

func InitUserPersistence(logger logger.Logger, db *persistencedb.PersistenceDB) storage.UserPersistence {
	return &userPersistence{
		logger: logger,
		db:     db,
	}
}

func (u *userPersistence) GetAllUsers(ctx context.Context, filters request_models.FilterParams) ([]dto.User, *model.MetaData, error) {
	users, total, err := u.db.GetAllUsersWithRole(ctx, utils.ComposeFilterSQL(ctx, filters, u.logger))
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no users found")
			u.logger.Info(ctx, "no users were found", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading users")
			u.logger.Error(ctx, "error reading users", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		}
	}
	return users, &model.MetaData{
		FilterParams: filters,
		Total:        total,
		Extra:        nil,
	}, nil
}

func (u *userPersistence) UpdateUserStatus(ctx context.Context, updateUserStatusParam dto.UpdateUserStatus, userID uuid.UUID) error {
	_, err := u.db.UpdateUser(ctx, db.UpdateUserParams{
		Status: sql.NullString{String: updateUserStatusParam.Status, Valid: true},
		ID:     userID,
	})

	if err != nil {
		err = errors.ErrUpdateError.Wrap(err, "error updating users")
		u.logger.Error(ctx, "error updating user's status", zap.Error(err), zap.Any("user-param", updateUserStatusParam))
		return err
	}
	return nil
}

func (u *userPersistence) UpdateUserRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	err := u.db.AssignRoleForUser(ctx, userID, roleName)
	if err != nil {
		err = errors.ErrUpdateError.Wrap(err, "error updating user role")
		u.logger.Error(ctx, "error updating user's role", zap.Error(err), zap.Any("user-id", userID), zap.String("role-name", roleName))
		return err
	}

	return nil
}

func (u *userPersistence) RevokeUserRole(ctx context.Context, userID uuid.UUID) error {
	err := u.db.RemoveRoleOFUser(ctx, userID)
	if err != nil {
		err = errors.ErrDBDelError.Wrap(err, "error revoking user role")
		u.logger.Error(ctx, "error revoking user's role", zap.Error(err), zap.Any("user-id", userID))
		return err
	}

	return nil
}

func (u *userPersistence) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	user, err := u.db.GetUserByIDWithRole(ctx, id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			u.logger.Info(ctx, "no user found", zap.Error(err), zap.String("id", id.String()))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			u.logger.Error(ctx, "unable to get user by id", zap.Error(err), zap.String("id", id.String()))
			return nil, err
		}
	}

	return user, nil
}
