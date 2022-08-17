package initiator

import (
	"sso/internal/storage"
	"sso/internal/storage/cache/authcode"
	"sso/internal/storage/cache/consent"
	"sso/internal/storage/cache/otp"
	"sso/internal/storage/cache/session"
	mock_otp "sso/mocks/storage/cache/otp"
	"sso/platform/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheLayer struct {
	OTPCacheLayer      storage.OTPCache
	SessionCacheLayer  storage.SessionCache
	ConsentCacheLayer  storage.ConsentCache
	AuthCodeCacheLayer storage.AuthCodeCache
}

type CacheOptions struct {
	OTPExpireTime      time.Duration
	SessionExpireTime  time.Duration
	ConsentExpireTime  time.Duration
	AuthCodeExpireTime time.Duration
}

func InitCacheLayer(client *redis.Client, options CacheOptions, log logger.Logger) CacheLayer {
	return CacheLayer{
		OTPCacheLayer:      otp.InitOTPCache(client, log.Named("otp-cache"), options.OTPExpireTime),
		SessionCacheLayer:  session.InitSessionCache(client, log.Named("session-cache"), options.SessionExpireTime),
		ConsentCacheLayer:  consent.InitConsentCache(client, log.Named("consent-cache"), options.ConsentExpireTime),
		AuthCodeCacheLayer: authcode.InitAuthCodeCache(client, log.Named("authcode-cache"), options.AuthCodeExpireTime),
	}
}

func InitMockCacheLayer(client *redis.Client, expireOn time.Duration, mockOTP string, log logger.Logger, options CacheOptions) CacheLayer {
	return CacheLayer{
		OTPCacheLayer:      mock_otp.InitMockOTPCache(client, log.Named("otp-cache"), options.OTPExpireTime, mockOTP),
		SessionCacheLayer:  session.InitSessionCache(client, log.Named("session-cache"), options.SessionExpireTime),
		ConsentCacheLayer:  consent.InitConsentCache(client, log.Named("consent-cache"), options.ConsentExpireTime),
		AuthCodeCacheLayer: authcode.InitAuthCodeCache(client, log.Named("authcode-cache"), options.AuthCodeExpireTime),
	}
}
