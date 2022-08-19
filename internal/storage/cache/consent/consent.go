package consent

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

type consentCache struct {
	logger   logger.Logger
	client   *redis.Client
	expireOn time.Duration
}

func InitConsentCache(client *redis.Client, log logger.Logger, exipreOn time.Duration) storage.ConsentCache {
	return &consentCache{log, client, exipreOn}
}

func (c *consentCache) GetConsent(ctx context.Context, consentID string) (dto.Consent, error) {
	consentKey := fmt.Sprintf(state.ConsentKey, consentID)
	consentResult, err := c.client.Get(ctx, consentKey).Result()
	if err != nil {
		if err == redis.Nil {
			err := errors.ErrCacheGetError.Wrap(err, "could not get from consent cache")
			c.logger.Error(ctx, "could not get from consent cache", zap.Error(err), zap.String("consentID", consentID))
			return dto.Consent{}, err
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from consent cache")
		c.logger.Error(ctx, "could not read from consent cache", zap.Error(err), zap.String("consentID", consentID))
		return dto.Consent{}, err
	}

	var consent dto.Consent
	err = json.Unmarshal([]byte(consentResult), &consent)
	if err != nil {
		err := errors.ErrCacheGetError.Wrap(err, "could not unmarshal consent")
		c.logger.Error(ctx, "could not unmarshal consent", zap.Error(err), zap.String("consentID", consentID))
		return dto.Consent{}, err
	}

	return consent, nil
}

func (c *consentCache) SaveConsent(ctx context.Context, consent dto.Consent) error {
	consentValue, err := json.Marshal(consent)
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not marshal consent")
		c.logger.Error(ctx, "could not marshal consent", zap.Error(err), zap.Any("consent", consent))
		return err
	}
	consentKey := fmt.Sprintf(state.ConsentKey, consent.ID)
	err = c.client.Set(ctx, consentKey, consentValue, c.expireOn).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set to consent cache")
		c.logger.Error(ctx, "could not write to consent cache", zap.Error(err), zap.Any("consent", consent))
	}
	return err
}

func (c *consentCache) DeleteConsent(ctx context.Context, consentID string) error {
	consentKey := fmt.Sprintf(state.ConsentKey, consentID)
	err := c.client.Del(ctx, consentKey).Err()
	if err != nil {
		err := errors.ErrCacheDel.Wrap(err, "could not delete from consent cache")
		c.logger.Error(ctx, "could not delete from consent cache", zap.Error(err), zap.String("consentID", consentID))
		return err
	}

	return nil
}

func (c *consentCache) ChangeStatus(ctx context.Context, status bool, consent dto.Consent) (dto.Consent, error) {
	consentKey := fmt.Sprintf(state.ConsentKey, consent.ID)
	consent.Approved = status
	consentValue, err := json.Marshal(consent)
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not marshal consent")
		c.logger.Error(ctx, "could not marshal consent", zap.Error(err), zap.Any("consent", consent))
		return dto.Consent{}, err
	}

	err = c.client.Set(ctx, consentKey, consentValue, redis.KeepTTL).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set to consent cache")
		c.logger.Error(ctx, "could not write to consent cache", zap.Error(err), zap.Any("consent", consent))
		return dto.Consent{}, err
	}

	return consent, nil
}
