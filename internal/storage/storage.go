package storage

import (
	"context"
	"sso/internal/constant/model/dto"
)

type OAuthPersistence interface {
	Register(ctx context.Context, user dto.User) (*dto.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*dto.User, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	UserByPhoneExists(ctx context.Context, phone string) (bool, error)
	UserByEmailExists(ctx context.Context, email string) (bool, error)
	GetUserByPhoneOrUserNameOrEmail(ctx context.Context, query string) (*dto.User, error)
}
