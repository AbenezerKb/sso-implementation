package initiator

import (
	"time"

	"sso/internal/storage"
	"sso/internal/storage/cache/authcode"
	"sso/internal/storage/cache/consent"
	"sso/internal/storage/cache/otp"
	"sso/internal/storage/cache/resetcode"
	"sso/internal/storage/cache/session"
	mock_otp "sso/mocks/storage/cache/otp"
	resetcode2 "sso/mocks/storage/cache/resetcode"
	"sso/platform/logger"

	"github.com/go-redis/redis/v8"
)

type CacheLayer struct {
	OTPCacheLayer       storage.OTPCache
	SessionCacheLayer   storage.SessionCache
	ConsentCacheLayer   storage.ConsentCache
	AuthCodeCacheLayer  storage.AuthCodeCache
	ResetCodeCacheLayer storage.ResetCodeCache
}

type CacheOptions struct {
	OTPExpireTime       time.Duration
	SessionExpireTime   time.Duration
	ConsentExpireTime   time.Duration
	AuthCodeExpireTime  time.Duration
	ResetCodeExpireTime time.Duration
}

func InitCacheLayer(client *redis.Client, options CacheOptions, log logger.Logger) CacheLayer {
	return CacheLayer{
		OTPCacheLayer:       otp.InitOTPCache(client, log.Named("otp-cache"), options.OTPExpireTime),
		SessionCacheLayer:   session.InitSessionCache(client, log.Named("session-cache"), options.SessionExpireTime),
		ConsentCacheLayer:   consent.InitConsentCache(client, log.Named("consent-cache"), options.ConsentExpireTime),
		AuthCodeCacheLayer:  authcode.InitAuthCodeCache(client, log.Named("authcode-cache"), options.AuthCodeExpireTime),
		ResetCodeCacheLayer: resetcode.InitResetCode(client, log.Named("reset-code-cache"), options.ResetCodeExpireTime),
	}
}

func InitMockCacheLayer(client *redis.Client, _ time.Duration, mockOTP string, log logger.Logger, options CacheOptions) CacheLayer {
	return CacheLayer{
		OTPCacheLayer:       mock_otp.InitMockOTPCache(client, log.Named("otp-cache"), options.OTPExpireTime, mockOTP),
		SessionCacheLayer:   session.InitSessionCache(client, log.Named("session-cache"), options.SessionExpireTime),
		ConsentCacheLayer:   consent.InitConsentCache(client, log.Named("consent-cache"), options.ConsentExpireTime),
		AuthCodeCacheLayer:  authcode.InitAuthCodeCache(client, log.Named("authcode-cache"), options.AuthCodeExpireTime),
		ResetCodeCacheLayer: resetcode2.InitMockResetCode(client, log.Named("reset-code-cache"), options.ResetCodeExpireTime, mockOTP),
	}
}
