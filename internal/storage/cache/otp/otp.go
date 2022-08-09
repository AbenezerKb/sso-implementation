package otp

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/storage"
	"sso/platform/logger"
	"time"
)

type otpCache struct {
	client   *redis.Client
	logger   logger.Logger
	expireOn time.Duration
}

func InitOTPCache(client *redis.Client, log logger.Logger, expireOn time.Duration) storage.OTPCache {
	return &otpCache{client, log, expireOn}
}

func (o *otpCache) GetOTP(ctx context.Context, phone string) (string, error) {
	otp, err := o.client.Get(ctx, phone).Result()
	if err != nil {
		if err == redis.Nil {
			return otp, nil
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from otp cache")
		o.logger.Error(ctx, "could not read from otp cache", zap.Error(err))
	}
	return otp, nil
}

func (o *otpCache) SetOTP(ctx context.Context, phone string, otp string) error {
	err := o.client.Set(ctx, phone, otp, o.expireOn).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set to otp cache")
		o.logger.Error(ctx, "could not set to otp cache", zap.Error(err))
		return err
	}

	return nil
}
