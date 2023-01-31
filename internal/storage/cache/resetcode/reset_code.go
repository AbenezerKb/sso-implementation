package resetcode

import (
	"context"
	"fmt"
	"time"

	"sso/internal/constant/errors"
	"sso/internal/constant/state"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ResetCode struct {
	logger   logger.Logger
	client   *redis.Client
	expireOn time.Duration
}

func InitResetCode(client *redis.Client, log logger.Logger, expireOn time.Duration) storage.ResetCodeCache {
	return &ResetCode{
		logger:   log,
		client:   client,
		expireOn: expireOn,
	}
}

func (c *ResetCode) GetResetCode(ctx context.Context, phone string) (string, error) {
	resetCodeKey := fmt.Sprintf(state.ResetCodeKey, phone)
	resetCodeResult, err := c.client.Get(ctx, resetCodeKey).Result()
	if err != nil {
		if err == redis.Nil {
			err := errors.ErrNoRecordFound.Wrap(err, "no record of reset code found")
			c.logger.Info(ctx, "reset code not found", zap.Error(err), zap.String("phone", phone))
			return "", err
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from reset code cache")
		c.logger.Error(ctx, "could not read from reset code cache", zap.Error(err))
		return "", err
	}

	return resetCodeResult, nil
}

func (c *ResetCode) SaveResetCode(ctx context.Context, phone, code string) error {
	resetCodeKey := fmt.Sprintf(state.ResetCodeKey, phone)
	err := c.client.Set(ctx, resetCodeKey, code, c.expireOn).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set reset code")
		c.logger.Error(ctx, "could not set reset code", zap.Error(err), zap.Any("phone", phone))
		return err
	}

	return nil
}

func (c *ResetCode) DeleteResetCode(ctx context.Context, phone string) error {
	resetCodeKey := fmt.Sprintf(state.ResetCodeKey, phone)
	err := c.client.Del(ctx, resetCodeKey).Err()
	if err != nil {
		err := errors.ErrCacheDel.Wrap(err, "could not delete reset code")
		c.logger.Error(ctx, "could not delete reset code", zap.Error(err), zap.String("phone", phone))
		return err
	}

	return nil
}
