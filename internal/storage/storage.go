package storage

import (
	"context"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"

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
	RemoveInternalRefreshToken(ctx context.Context, refreshToken string) error
	SaveInternalRefreshToken(ctx context.Context, rf dto.InternalRefreshToken) error
	GetInternalRefreshToken(ctx context.Context, refreshtoken string) (*dto.InternalRefreshToken, error)
	UpdateInternalRefreshToken(ctx context.Context, param dto.InternalRefreshToken) (*dto.InternalRefreshToken, error)
	GetInternalRefreshTokenByUserID(ctx context.Context, userID uuid.UUID) (*dto.InternalRefreshToken, error)
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
	GetNamedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error)
	AuthHistoryExists(ctx context.Context, code string) (bool, error)
	PersistRefreshToken(ctx context.Context, param dto.RefreshToken) (*dto.RefreshToken, error)
	RemoveRefreshTokenCode(ctx context.Context, code string) error
	RemoveRefreshToken(ctx context.Context, refresh_token string) error
	AddAuthHistory(ctx context.Context, param dto.AuthHistory) (*dto.AuthHistory, error)
	CheckIfUserGrantedClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) (bool, dto.RefreshToken, error)
	GetRefreshToken(ctx context.Context, token string) (*dto.RefreshToken, error)
	GetRefreshTokenOfClientByUserID(ctx context.Context, userID, clientID uuid.UUID) (*dto.RefreshToken, error)
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
	DeleteClientByID(ctx context.Context, id uuid.UUID) error
	GetAllClients(ctx context.Context, filters request_models.FilterParams) ([]dto.Client, *model.MetaData, error)
}

type AuthCodeCache interface {
	SaveAuthCode(ctx context.Context, authCode dto.AuthCode) error
	GetAuthCode(ctx context.Context, code string) (dto.AuthCode, error)
	DeleteAuthCode(ctx context.Context, code string) error
}

type ScopePersistence interface {
	CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error)
	GetScope(ctx context.Context, scope string) (dto.Scope, error)
	GetListedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error)
	GetScopeNameOnly(ctx context.Context, scopes ...string) (string, error)
}

type UserPersistence interface {
}

type ProfilePersistence interface {
	UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error)
}
