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
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userPersistence struct {
	logger logger.Logger
	db     *db.Queries
}

func InitUserPersistence(logger logger.Logger, db *db.Queries) storage.UserPersistence {
	return &userPersistence{
		logger: logger,
		db:     db,
	}
}

func (u *userPersistence) GetAllUsers(ctx context.Context, filters request_models.FilterParams) ([]dto.User, *model.MetaData, error) {
	users, total, err := u.db.GetAllUsers(ctx, utils.ComposeFilterSQL(ctx, filters, u.logger))
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
	usersDTO := make([]dto.User, len(users))
	for k, v := range users {
		usersDTO[k] = dto.User{
			ID:         v.ID,
			Status:     v.Status.String,
			FirstName:  v.FirstName,
			MiddleName: v.MiddleName,
			LastName:   v.LastName,
			Email:      v.Email.String,
			Phone:      v.Phone,
			Gender:     v.Gender,
			CreatedAt:  v.CreatedAt,
		}
	}
	return usersDTO, &model.MetaData{
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
