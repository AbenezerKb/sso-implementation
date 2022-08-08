package module

import (
	"context"
	"sso/internal/constant/model/dto"
)

type OAuthModule interface {
	Register(ctx context.Context, user dto.User) (*dto.User, error)
	Login(ctx context.Context, user dto.User) (*dto.TokenResponse, error)
	ComparePassword(hashedPwd, plainPassword string) bool
	HashAndSalt(pwd []byte) string
}
