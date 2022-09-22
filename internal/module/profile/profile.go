package profile

import (
	"context"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type profileModule struct {
	logger             logger.Logger
	oauthPersistence   storage.OAuthPersistence
	profilePersistence storage.ProfilePersistence
}

func InitProfile(logger logger.Logger, oauthPersistence storage.OAuthPersistence, profilePersistence storage.ProfilePersistence) module.ProfileModule {
	return &profileModule{
		logger:             logger,
		oauthPersistence:   oauthPersistence,
		profilePersistence: profilePersistence,
	}
}

func (p *profileModule) UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error) {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return &dto.User{}, err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return nil, err
	}

	if err := userParam.ValidateUpdateProfile(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	userParam.ID = userID
	updatedUser, err := p.profilePersistence.UpdateProfile(ctx, userParam)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
