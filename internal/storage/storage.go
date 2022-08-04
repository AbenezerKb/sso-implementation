package storage

import (
	"context"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
)

type OAuthPersistence interface {
	Register(ctx context.Context, user dto.User) (*db.User, error)
	GetUserByPhone(ctx context.Context, phone string) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	UserByPhoneExists(ctx context.Context, phone string) (bool, error)
	UserByEmailExists(ctx context.Context, email string) (bool, error)
}
