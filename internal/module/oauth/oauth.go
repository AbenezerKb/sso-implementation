package oauth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/go-redis/redis/v8"
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

func (o *oauth) Register(ctx context.Context, userParam dto.User) (*dto.User, error) {
	if err := userParam.ValidateUser(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	exists, err := o.oauthPersistence.UserByPhoneExists(ctx, userParam.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrDataExists.Wrap(err, "user with this phone already exists")
	}

	if userParam.Email != "" {
		exists, err := o.oauthPersistence.UserByEmailExists(ctx, userParam.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrDataExists.Wrap(err, "user with this email already exists")
		}
	}

	userParam.Password, err = o.HashAndSalt(ctx, []byte(userParam.Password))
	if err != nil {
		return nil, err
	}

	user, err := o.oauthPersistence.Register(ctx, userParam)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (o *oauth) HashAndSalt(ctx context.Context, pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		o.logger.Error(ctx, "could not hash password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

func (o *oauth) ComparePassword(hashedPwd, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPassword))
	return err == nil
}
