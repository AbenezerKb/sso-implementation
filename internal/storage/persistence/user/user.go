package user

import (
	"context"
	"database/sql"

	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
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

func (u *userPersistence) GetAllUsers(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.User, *model.MetaData, error) {
	users, total, err := u.db.GetAllUsersWithRole(ctx, filters)
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

func (u *userPersistence) GetUserByPhone(ctx context.Context, phone string) (*dto.User, error) {
	user, err := u.db.GetUserByPhone(ctx, phone)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			u.logger.Info(ctx, "no user found", zap.Error(err), zap.String("phone", phone))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			u.logger.Error(ctx, "unable to get user by phone", zap.Error(err), zap.String("phone", phone))
			return nil, err
		}
	}

	return &dto.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		MiddleName:     user.MiddleName,
		LastName:       user.LastName,
		Email:          user.Email.String,
		Phone:          user.Phone,
		UserName:       user.UserName,
		Gender:         user.Gender,
		Status:         user.Status.String,
		ProfilePicture: user.ProfilePicture.String,
		CreatedAt:      user.CreatedAt,
	}, nil
}
func (u *userPersistence) GetUsersByPhone(ctx context.Context, phones []string) ([]dto.User, error) {
	users, err := u.db.GetUsersByParsedField(ctx, "phone", phones)
	if err != nil {
		err := errors.ErrReadError.Wrap(err, "error fetching users")
		u.logger.Error(ctx, "error while fetching users by phone number")
		return nil, err
	}

	return users, nil
}
func (u *userPersistence) GetUsersByID(ctx context.Context, ids []string) ([]dto.User, error) {
	users, err := u.db.GetUsersByParsedField(ctx, "id", ids)
	if err != nil {
		err := errors.ErrReadError.Wrap(err, "error fetching users")
		u.logger.Error(ctx, "error while fetching users by id")
		return nil, err
	}

	return users, nil
}

func (u *userPersistence) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) (*dto.User, error) {
	user, err := u.db.ChangeUserPasswordByID(ctx, db.ChangeUserPasswordByIDParams{
		Password: newPassword,
		ID:       userID,
	})
	if err != nil {
		err := errors.ErrUpdateError.Wrap(err, "unable to update user password")
		u.logger.Error(ctx, "error while trying to update user password",
			zap.Error(err),
			zap.Any("user-id", userID))

		return nil, err
	}

	return &dto.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		MiddleName:     user.MiddleName,
		LastName:       user.LastName,
		Email:          user.Email.String,
		Phone:          user.Phone,
		UserName:       user.UserName,
		Gender:         user.Gender,
		Status:         user.Status.String,
		ProfilePicture: user.ProfilePicture.String,
		CreatedAt:      user.CreatedAt,
	}, nil
}
func (u *userPersistence) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := u.db.RemoveUser(ctx, userID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "user not found")
			u.logger.Info(ctx, "user was not found", zap.Error(err),
				zap.String("user-ID", userID.String()))
			return err
		}
		err := errors.ErrDBDelError.Wrap(err, "error deleting user")
		u.logger.Error(ctx, "error while deleting user", zap.Error(err),
			zap.String("user-ID", userID.String()))
		return err
	}
	return nil
}
