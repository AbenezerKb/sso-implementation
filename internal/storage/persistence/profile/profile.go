package profile

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"

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
