package module

import (
	"context"
	"crypto/rsa"
	"sso/internal/constant/model/dto"

	"github.com/golang-jwt/jwt/v4"
)

type OAuthModule interface {
	Register(ctx context.Context, user dto.RegisterUser) (*dto.User, error)
	Login(ctx context.Context, login dto.LoginCredential) (*dto.TokenResponse, error)
	ComparePassword(hashedPwd, plainPassword string) bool
	HashAndSalt(ctx context.Context, pwd []byte) (string, error)
	RequestOTP(ctx context.Context, phone string, rqType string) error
	VerifyToken(signingMethod jwt.SigningMethod, token string, pk *rsa.PublicKey) (bool, *jwt.RegisteredClaims)
	GetUserStatus(ctx context.Context, Id string) (string, error)
}

type UserModule interface {
	Create(ctx context.Context, user dto.CreateUser) (*dto.User, error)
}

type ClientModule interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
}
