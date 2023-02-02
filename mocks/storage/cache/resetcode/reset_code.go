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

type resetCode struct {
	logger   logger.Logger
	client   *redis.Client
	expireOn time.Duration
	mockCode string
}

func InitMockResetCode(client *redis.Client, log logger.Logger, expireOn time.Duration, mockCode string) storage.ResetCodeCache {
	return &resetCode{
		logger:   log,
		client:   client,
		expireOn: expireOn,
		mockCode: mockCode,
	}
}

func (c *resetCode) GetResetCode(ctx context.Context, email string) (string, error) {
	resetCodeKey := fmt.Sprintf(state.ResetCodeKey, email)
	resetCodeResult, err := c.client.Get(ctx, resetCodeKey).Result()
	if err != nil {
		if err == redis.Nil {
			err := errors.ErrNoRecordFound.Wrap(err, "no record of reset code found")
			c.logger.Info(ctx, "reset code not found", zap.Error(err), zap.String("email", email))
			return "", err
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from reset code cache")
		c.logger.Error(ctx, "could not read from reset code cache", zap.Error(err))
		return "", err
	}

	return resetCodeResult, nil
}

func (c *resetCode) SaveResetCode(ctx context.Context, email, _ string) error {
	resetCodeKey := fmt.Sprintf(state.ResetCodeKey, email)
	err := c.client.Set(ctx, resetCodeKey, c.mockCode, c.expireOn).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set reset code")
		c.logger.Error(ctx, "could not set reset code", zap.Error(err), zap.Any("email", email))
		return err
	}

	return nil
}

func (c *resetCode) DeleteResetCode(ctx context.Context, email string) error {
	resetCodeKey := fmt.Sprintf(state.ResetCodeKey, email)
	err := c.client.Del(ctx, resetCodeKey).Err()
	if err != nil {
		err := errors.ErrCacheDel.Wrap(err, "could not delete reset code")
		c.logger.Error(ctx, "could not delete reset code", zap.Error(err), zap.String("email", email))
		return err
	}

	return nil
}
