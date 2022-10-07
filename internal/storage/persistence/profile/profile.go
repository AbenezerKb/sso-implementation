package profile

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type profilePersistence struct {
	logger logger.Logger
	db     *db.Queries
}

func InitProfilePersistence(logger logger.Logger, db *db.Queries) storage.ProfilePersistence {
	return &profilePersistence{
		logger: logger,
		db:     db,
	}
}

func (p *profilePersistence) UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error) {
	user, err := p.db.UpdateProfile(ctx, db.UpdateProfileParams{
		FirstName:  userParam.FirstName,
		MiddleName: userParam.MiddleName,
		LastName:   userParam.LastName,
		Gender:     userParam.Gender,
		ID:         userParam.ID,
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not update user profile")
		p.logger.Error(ctx, "unable to update user profile", zap.Error(err), zap.Any("user", userParam))
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
		Gender:         user.Gender,
		ProfilePicture: user.ProfilePicture.String,
	}, nil
}

func (p *profilePersistence) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.User, error) {
	user, err := p.db.GetUserById(ctx, userID)
	if err != nil {
		err = errors.ErrReadError.Wrap(err, "could not reade user profile")
		p.logger.Error(ctx, "unable to reade user profile", zap.Error(err), zap.Any("user-id", userID))
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
		ProfilePicture: user.ProfilePicture.String,
		CreatedAt:      user.CreatedAt,
	}, nil
}

func (p *profilePersistence) UpdateProfilePicture(ctx context.Context, finalImageName string, userID uuid.UUID) error {
	_, err := p.db.UpdateUser(ctx, db.UpdateUserParams{
		ProfilePicture: sql.NullString{String: finalImageName, Valid: true},
		ID:             userID,
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not update user profile picture")
		p.logger.Error(ctx, "unable to update user profile picture", zap.Error(err), zap.Any("imageName", finalImageName), zap.Any("user-id", userID))
		return err
	}

	return nil
}

func (p *profilePersistence) ChangePhone(ctx context.Context, changePhoneParam dto.ChangePhoneParam, userID uuid.UUID) error {
	_, err := p.db.UpdateUser(ctx, db.UpdateUserParams{
		Phone: sql.NullString{String: changePhoneParam.Phone, Valid: true},
		ID:    userID,
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not change user phone number")
		p.logger.Error(ctx, "unable to update user's phone number", zap.Error(err), zap.Any("phone", changePhoneParam.Phone), zap.Any("user-id", userID))
		return err
	}

	return nil

}

func (p *profilePersistence) ChangePassword(ctx context.Context, changePasswordParam dto.ChangePasswordParam, userID uuid.UUID) error {
	_, err := p.db.UpdateUser(ctx, db.UpdateUserParams{
		Password: sql.NullString{String: changePasswordParam.NewPassword, Valid: true},
		ID:       userID,
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not change user password")
		p.logger.Error(ctx, "unable to change user's password", zap.Error(err), zap.Any("user-id", userID))
		return err
	}

	return nil
}
