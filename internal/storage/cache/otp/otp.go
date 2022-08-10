package otp

import (
	"context"
	"fmt"
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
func (o *otpCache) GetDelOTP(ctx context.Context, phone string) (string, error) {
	otp, err := o.client.GetDel(ctx, phone).Result()
	if err != nil {
		if err == redis.Nil {
			return otp, nil
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from otp cache")
		o.logger.Error(ctx, "could not read from otp cache", zap.Error(err))
	}
	return otp, nil
}

func (o *otpCache) DeleteOTP(ctx context.Context, phone ...string) error {
	err := o.client.Del(ctx, phone...).Err()
	if err != nil {
		err := errors.ErrCacheDel.Wrap(err, fmt.Sprintf("couldn't delete cache"))
		o.logger.Error(ctx, fmt.Sprintf("couldn't delete caches: %v", phone), zap.Error(err))
		return err
	}

	return nil
}
