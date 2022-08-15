package oauth2

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type oauth2 struct {
	logger logger.Logger
	db     *db.Queries
}

func InitOAuth2(logger logger.Logger, db *db.Queries) storage.OAuth2Persistence {
	return &oauth2{
		logger,
		db,
	}
}

func (o *oauth2) GetClient(ctx context.Context, id string) (*dto.Client, error) {
	return &dto.Client{
		RedirectURIs: []string{"http://localhost:9000/callback"},
		Scope:        "openid profile email",
		Name:         "test",
		Secret:       "test",
		ID:           "test",
	}, nil
}

func (o *oauth2) GetNamedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error) {
	return []dto.Scope{
		{Name: "openid", Description: "openid"},
		{Name: "profile", Description: "profile"},
		{Name: "email", Description: "email"},
	}, nil
}

func (o *oauth2) SaveAuthCode(ctx context.Context, authCode dto.AuthCode) error {
	client_id, err := uuid.Parse(authCode.ClientID)
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not parse client id")
		o.logger.Error(ctx, zap.Error(err).String)
		return err
	}
	user_id, err := uuid.Parse(authCode.UserID)
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not parse user id")
		o.logger.Error(ctx, zap.Error(err).String)
		return err
	}

	o.db.CreateAuthCode(ctx, db.CreateAuthCodeParams{
		Code:        authCode.Code,
		ClientID:    client_id,
		UserID:      user_id,
		RedirectUri: authCode.RedirectURI,
		Scope:       authCode.Scope,
	})
	return nil
}

func (o *oauth2) GetAuthCode(ctx context.Context, code string) (dto.AuthCode, error) {
	authCode, err := o.db.GetAuthCode(ctx, code)
	if err != nil {
		err = errors.ErrReadError.Wrap(err, "could not read auth code")
		o.logger.Error(ctx, zap.Error(err).String)
		return dto.AuthCode{}, err
	}
	return dto.AuthCode{
		Code:        authCode.Code,
		ClientID:    authCode.ClientID.String(),
		UserID:      authCode.UserID.String(),
		RedirectURI: authCode.RedirectUri,
		Scope:       authCode.Scope,
	}, nil
}

func (o *oauth2) DeleteAuthCode(ctx context.Context, code string) error {
	_, err := o.db.DeleteAuthCode(ctx, code)
	if err != nil {
		err = errors.ErrNoRecordFound.Wrap(err, "could not delete auth code")
		o.logger.Error(ctx, zap.Error(err).String)
		return err
	}

	return nil
}
