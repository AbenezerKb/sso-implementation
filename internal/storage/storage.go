package storage

import (
	"context"
	"sso/internal/constant/model/dto"

	"github.com/google/uuid"
)

type OAuthPersistence interface {
	Register(ctx context.Context, user dto.User) (*dto.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*dto.User, error)
	GetUserStatus(ctx context.Context, Id uuid.UUID) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	UserByPhoneExists(ctx context.Context, phone string) (bool, error)
	UserByEmailExists(ctx context.Context, email string) (bool, error)
	GetUserByPhoneOrEmail(ctx context.Context, query string) (*dto.User, error)
}

type OTPCache interface {
	SetOTP(ctx context.Context, phone string, otp string) error
	GetOTP(ctx context.Context, phone string) (string, error)
	GetDelOTP(ctx context.Context, phone string) (string, error)
	DeleteOTP(ctx context.Context, phone ...string) error
}

type SessionCache interface {
	SaveSession(ctx context.Context, session dto.Session) error
	GetSession(ctx context.Context, sessionID string) (dto.Session, error)
}

type ClientPersistence interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
}
