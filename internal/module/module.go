package module

import (
	"context"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"

	"github.com/joomcode/errorx"

	"github.com/google/uuid"
)

type OAuthModule interface {
	Register(ctx context.Context, user dto.RegisterUser) (*dto.User, error)
	Login(ctx context.Context, login dto.LoginCredential) (*dto.TokenResponse, error)
	ComparePassword(hashedPwd, plainPassword string) bool
	RequestOTP(ctx context.Context, phone string, rqType string) error
	GetUserStatus(ctx context.Context, Id string) (string, error)
	Logout(ctx context.Context, param dto.InternalRefreshTokenRequestBody) error
	RefreshToken(ctx context.Context, param dto.InternalRefreshTokenRequestBody) (*dto.TokenResponse, error)
}

type OAuth2Module interface {
	Authorize(ctx context.Context, authRequestParma dto.AuthorizationRequestParam, requestOrigin string, bindError *errorx.Error) string
	GetConsentByID(ctx context.Context, consentID string) (dto.ConsentResponse, error)
	ApproveConsent(ctx context.Context, consentID string, userID uuid.UUID, opbs string, bindError *errorx.Error) string
	RejectConsent(ctx context.Context, consentID, failureReason string, bindError *errorx.Error) string
	Token(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error)
	Logout(ctx context.Context, logoutReqParam dto.LogoutRequest, bindError *errorx.Error) string
	RevokeClient(ctx context.Context, clientBody request_models.RevokeClientBody) error
}
type UserModule interface {
	Create(ctx context.Context, user dto.CreateUser) (*dto.User, error)
	UpdateProfile(ctx context.Context, user dto.User) (*dto.User, error)
	GetUserByID(ctx context.Context, id string) (*dto.User, error)
}

type ClientModule interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*dto.Client, error)
	DeleteClientByID(ctx context.Context, id string) error
}

type ScopeMoudle interface {
	GetScope(ctx context.Context, scope string) (dto.Scope, error)
	CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error)
}
