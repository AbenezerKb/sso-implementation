package module

import (
	"context"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
)

type OAuthModule interface {
	Register(ctx context.Context, user dto.User) (*db.User, error)
}
