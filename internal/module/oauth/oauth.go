package oauth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
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
		o.logger.Error(ctx, "invalid input", zap.Error(err))
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

	userParam.Password = o.HashAndSalt([]byte(userParam.Password))
	user, err := o.oauthPersistence.Register(ctx, userParam)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (o *oauth) HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
func (o *oauth) ComparePassword(hashedPwd, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPassword))
	return err == nil
}
