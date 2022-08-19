package oauth2

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

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

func (o *oauth2) GetClient(ctx context.Context, id uuid.UUID) (*dto.Client, error) {
	return &dto.Client{
		RedirectURIs: []string{"https://www.google.com/"},
		Scopes:       "openid profile email",
		Name:         "test",
		Secret:       "test",
		// ID:           "test",
	}, nil
}

func (o *oauth2) GetNamedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error) {
	return []dto.Scope{
		{Name: "openid", Description: "openid"},
		{Name: "profile", Description: "profile"},
		{Name: "email", Description: "email"},
	}, nil
}

func (o *oauth2) AuthHistoryExists(ctx context.Context, code string) (bool, error) {
	_, err := o.db.GetAuthHistory(ctx, code)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		}
		err = errors.ErrReadError.Wrap(err, "could not read auth history")
		o.logger.Error(ctx, "error reading from auth history", zap.Error(err), zap.String("code", code))
		return false, err
	}
	return true, nil
}

func (o *oauth2) PersistRefreshToken(ctx context.Context, param dto.RefreshToken) (*dto.RefreshToken, error) {
	refToken, err := o.db.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
		ExpiresAt:    param.ExpiresAt,
		UserID:       param.UserID,
		ClientID:     param.ClientID,
		Scope:        utils.StringOrNull(param.Scope),
		RedirectUri:  utils.StringOrNull(param.RedirectUri),
		Refreshtoken: param.Refreshtoken,
		Code:         param.Code,
	})
	if err != nil {
		Err := errors.ErrWriteError.Wrap(err, "unable to persist the refresh token")
		o.logger.Error(ctx, "error saving the refresh token", zap.Error(Err), zap.Any("refresh-token", param))
		return nil, Err
	}
	return &dto.RefreshToken{
		Code:         refToken.Code,
		Refreshtoken: refToken.Refreshtoken,
		RedirectUri:  refToken.RedirectUri.String,
		Scope:        refToken.Scope.String,
		UserID:       refToken.UserID,
		ID:           refToken.ID,
		ClientID:     refToken.ClientID,
	}, nil
}

func (o *oauth2) AddAuthHistory(ctx context.Context, param dto.AuthHistory) (*dto.AuthHistory, error) {
	authHist, err := o.db.CreateAuthHistory(ctx, db.CreateAuthHistoryParams{
		UserID:      param.UserID,
		ClientID:    param.ClientID,
		Scope:       utils.StringOrNull(param.Scope),
		RedirectUri: utils.StringOrNull(param.RedirectUri),
		Status:      param.Status,
		Code:        param.Code,
	})
	if err != nil {
		Err := errors.ErrWriteError.Wrap(err, " unable to save the auth history")
		o.logger.Error(ctx, "error saving the auth history", zap.Error(Err), zap.Any("auth-history", param))
		return nil, Err
	}
	return &dto.AuthHistory{
		ID:          authHist.ID,
		UserID:      authHist.UserID,
		ClientID:    authHist.ClientID,
		RedirectUri: authHist.RedirectUri.String,
		Scope:       authHist.Scope.String,
		Code:        authHist.Code,
		Status:      authHist.Status,
	}, nil
}

func (o *oauth2) RemoveRefreshToken(ctx context.Context, code string) error {
	if err := o.db.RemoveRefreshToken(ctx, code); err != nil {
		err := errors.ErrDBDelError.Wrap(err, "could be able to delete the referesh token")
		o.logger.Error(ctx, "unable to delete the refresh token", zap.Error(err), zap.Any("refresh-token-code", code))
		return err
	}
	return nil
}
