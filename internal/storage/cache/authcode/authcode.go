package authcode

import (
	"context"
	"encoding/json"
	"fmt"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/state"
	"sso/internal/storage"
	"sso/platform/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type AuthCode struct {
	logger   logger.Logger
	client   *redis.Client
	expireOn time.Duration
}

func InitAuthCodeCache(client *redis.Client, log logger.Logger, exipreOn time.Duration) storage.AuthCodeCache {
	return &AuthCode{
		logger:   log,
		client:   client,
		expireOn: exipreOn,
	}
}

func (c *AuthCode) GetAuthCode(ctx context.Context, code string) (dto.AuthCode, error) {
	authCodeKey := fmt.Sprintf(state.AuthCodeKey, code)
	authCodeResult, err := c.client.Get(ctx, authCodeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return dto.AuthCode{}, nil
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from authcode cache")
		c.logger.Error(ctx, "could not read from authcode cache", zap.Error(err))
	}

	var authCode dto.AuthCode
	err = json.Unmarshal([]byte(authCodeResult), &authCode)
	if err != nil {
		err := errors.ErrCacheGetError.Wrap(err, "could not unmarshal authcode")
		c.logger.Error(ctx, "could not unmarshal authcode", zap.Error(err))
		return dto.AuthCode{}, err
	}

	return authCode, nil
}

func (c *AuthCode) SaveAuthCode(ctx context.Context, authCode dto.AuthCode) error {
	authCodeValue, err := json.Marshal(authCode)
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not marshal authcode")
		c.logger.Error(ctx, "could not marshal authcode", zap.Error(err))
		return err
	}
	authCodeKey := fmt.Sprintf(state.AuthCodeKey, authCode.Code)
	err = c.client.Set(ctx, authCodeKey, authCodeValue, c.expireOn).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set authcode")
		c.logger.Error(ctx, "could not set authcode", zap.Error(err))
		return err
	}

	return nil
}

func (c *AuthCode) DeleteAuthCode(ctx context.Context, code string) error {
	authCodeKey := fmt.Sprintf(state.AuthCodeKey, code)
	err := c.client.Del(ctx, authCodeKey).Err()
	if err != nil {
		err := errors.ErrCacheDel.Wrap(err, "could not delete authcode")
		c.logger.Error(ctx, "could not delete authcode", zap.Error(err))
		return err
	}

	return nil
}
