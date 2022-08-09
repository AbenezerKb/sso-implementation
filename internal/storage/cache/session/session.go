package session

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
	"time"
)

type sessionCache struct {
	client   *redis.Client
	logger   logger.Logger
	expireOn time.Duration
}

func InitSessionCache(client *redis.Client, log logger.Logger, expireOn time.Duration) storage.SessionCache {
	return &sessionCache{client, log, expireOn}
}

func (s *sessionCache) SaveSession(ctx context.Context, session dto.Session) error {
	sessionValue, err := json.Marshal(session)
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not marshal session")
		s.logger.Error(ctx, "could not marshal session", zap.Error(err))
		return err
	}
	err = s.client.Set(ctx, session.ID, sessionValue, s.expireOn).Err()
	if err != nil {
		err := errors.ErrCacheSetError.Wrap(err, "could not set to session cache")
		s.logger.Error(ctx, "could not set to session cache", zap.Error(err))
		return err
	}

	return nil
}

func (s *sessionCache) GetSession(ctx context.Context, sessionID string) (dto.Session, error) {
	sessionValue, err := s.client.Get(ctx, sessionID).Result()
	if err != nil {
		if err == redis.Nil {
			return dto.Session{}, nil
		}

		err := errors.ErrCacheGetError.Wrap(err, "could not get from session cache")
		s.logger.Error(ctx, "could not read from session cache", zap.Error(err))
	}
	var session dto.Session
	err = json.Unmarshal([]byte(sessionValue), &session)
	if err != nil {
		err := errors.ErrCacheGetError.Wrap(err, "could not unmarshal session")
		s.logger.Error(ctx, "could not unmarshal session", zap.Error(err))
		return dto.Session{}, err
	}
	return session, nil
}
