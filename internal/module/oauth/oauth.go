package oauth

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/go-redis/redis/v8"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type oauth struct {
	logger           logger.Logger
	oauthPersistence storage.OAuthPersistence
	cache            *redis.Client
}

func InitOAuth(logger logger.Logger, oauthPersistence storage.OAuthPersistence, cache *redis.Client) module.OAuthModule {
	return &oauth{
		logger,
		oauthPersistence,
		cache,
	}
}

func (o *oauth) Register(ctx context.Context, userParam dto.User) (*db.User, error) {
	if err := userParam.ValidateUser(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Error(ctx, zap.Error(err).String)
		return nil, err
	}

	_, err := o.oauthPersistence.GetUserByPhone(ctx, userParam.Phone)
	if err != nil && !errorx.IsOfType(err, errors.ErrNoRecordFound) {
		return nil, err
	}

	if userParam.Email != "" {
		_, err := o.oauthPersistence.GetUserByEmail(ctx, userParam.Email)
		if err != nil {
			return nil, err
		}
	}

	user, err := o.oauthPersistence.Register(ctx, userParam)
	if err != nil {
		return nil, err
	}
	return user, nil
}
