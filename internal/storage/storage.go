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
	GetUserByID(ctx context.Context, Id uuid.UUID) (*dto.User, error)
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

type OAuth2Persistence interface {
	GetClient(ctx context.Context, id uuid.UUID) (*dto.Client, error)
	GetNamedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error)
	AuthHistoryExists(ctx context.Context, code string) (bool, error)
	PersistRefreshToken(ctx context.Context, param dto.RefreshToken) (*dto.RefreshToken, error)
	RemoveRefreshToken(ctx context.Context, code string) error
	AddAuthHistory(ctx context.Context, param dto.AuthHistory) (*dto.AuthHistory, error)
	CheckIfUserGrantedClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) (bool, dto.RefreshToken, error)
}

type ConsentCache interface {
	SaveConsent(ctx context.Context, consent dto.Consent) error
	GetConsent(ctx context.Context, consentID string) (dto.Consent, error)
	DeleteConsent(ctx context.Context, consentID string) error
	ChangeStatus(ctx context.Context, status bool, consent dto.Consent) (dto.Consent, error)
}
type ClientPersistence interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*dto.Client, error)
}

type AuthCodeCache interface {
	SaveAuthCode(ctx context.Context, authCode dto.AuthCode) error
	GetAuthCode(ctx context.Context, code string) (dto.AuthCode, error)
	DeleteAuthCode(ctx context.Context, code string) error
}
