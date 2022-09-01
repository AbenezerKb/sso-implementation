package user

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"

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

func (u *userPersistence) UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error) {
	updateData := db.UpdateProfileParams{}

	if userParam.Email != "" {
		updateData.Email = sql.NullString{String: userParam.Email, Valid: true}
	}

	if userParam.FirstName != "" {
		updateData.FirstName = sql.NullString{String: userParam.FirstName, Valid: true}
	}

	if userParam.MiddleName != "" {
		updateData.MiddleName = sql.NullString{String: userParam.MiddleName, Valid: true}
	}

	if userParam.LastName != "" {
		updateData.LastName = sql.NullString{String: userParam.LastName, Valid: true}
	}

	if userParam.Phone != "" {
		updateData.Phone = sql.NullString{String: userParam.Phone, Valid: true}
	}

	if userParam.Gender != "" {
		updateData.Gender = sql.NullString{String: userParam.Gender, Valid: true}
	}

	if userParam.UserName != "" {
		updateData.UserName = sql.NullString{String: userParam.UserName, Valid: true}
	}

	updateData.ID = userParam.ID
	user, err := u.db.UpdateProfile(ctx, updateData)

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not update user profile")
		u.logger.Error(ctx, "unable to update user profile", zap.Error(err), zap.Any("user", userParam))
		return &dto.User{}, err
	}

	return &dto.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		MiddleName:     user.MiddleName,
		LastName:       user.MiddleName,
		Email:          user.Email.String,
		Phone:          user.Phone,
		UserName:       user.UserName,
		ProfilePicture: user.ProfilePicture.String,
	}, nil
}
