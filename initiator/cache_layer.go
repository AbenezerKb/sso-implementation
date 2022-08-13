package initiator

import (
	"sso/internal/storage"
	"sso/internal/storage/cache/consent"
	"sso/internal/storage/cache/otp"
	"sso/internal/storage/cache/session"
	mock_otp "sso/mocks/storage/cache/otp"
	"sso/platform/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheLayer struct {
	OTPCacheLayer     storage.OTPCache
	SessionCacheLayer storage.SessionCache
	ConsentCacheLayer storage.ConsentCache
}

func InitCacheLayer(client *redis.Client, expireOn time.Duration, log logger.Logger) CacheLayer {
	return CacheLayer{
		OTPCacheLayer:     otp.InitOTPCache(client, log, expireOn),
		SessionCacheLayer: session.InitSessionCache(client, log, expireOn),
		ConsentCacheLayer: consent.InitConsentCache(client, log, expireOn),
	}
}

func InitMockCacheLayer(client *redis.Client, expireOn time.Duration, mockOTP string, log logger.Logger) CacheLayer {
	return CacheLayer{
		OTPCacheLayer:     mock_otp.InitMockOTPCache(client, log, expireOn, mockOTP),
		SessionCacheLayer: session.InitSessionCache(client, log, expireOn),
	}
}
