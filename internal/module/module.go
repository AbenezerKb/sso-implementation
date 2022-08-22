package module

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"

	"github.com/google/uuid"
)

type OAuthModule interface {
	Register(ctx context.Context, user dto.RegisterUser) (*dto.User, error)
	Login(ctx context.Context, login dto.LoginCredential) (*dto.TokenResponse, error)
	ComparePassword(hashedPwd, plainPassword string) bool
	RequestOTP(ctx context.Context, phone string, rqType string) error
	GetUserStatus(ctx context.Context, Id string) (string, error)
}

type OAuth2Module interface {
	Authorize(ctx context.Context, authRequestParma dto.AuthorizationRequestParam) (string, errors.AuhtErrResponse, error)
	GetConsentByID(ctx context.Context, consentID string, id string) (dto.ConsentData, error)
	ApproveConsent(ctx context.Context, consentID string, userID uuid.UUID) (string, error)
	RejectConsent(ctx context.Context, consentID, failureReason string) (string, error)
	IssueAuthCode(ctx context.Context, consent dto.Consent) (string, string, error)
	Token(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error)
}
type UserModule interface {
	Create(ctx context.Context, user dto.CreateUser) (*dto.User, error)
}

type ClientModule interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*dto.Client, error)
}

type ScopeMoudle interface {
	GetScope(ctx context.Context, scope string) (dto.Scope, error)
	CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error)
}
