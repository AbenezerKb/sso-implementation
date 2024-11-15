package oauth2

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"
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

func (o *oauth2) GetNamedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error) {
	namedScopes := []dto.Scope{}
	for _, scope := range scopes {
		scope, err := o.db.GetScope(ctx, scope)
		if err != nil {
			if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
				continue
			}
			err = errors.ErrReadError.Wrap(err, "could not read the scope")
			o.logger.Error(ctx, "unable to read the scope", zap.Error(err), zap.Any("scope", scope))
			return nil, err
		}
		namedScopes = append(namedScopes, dto.Scope{
			Name:        scope.Name,
			Description: scope.Description,
		})

	}
	return namedScopes, nil
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
		RefreshToken: param.RefreshToken,
		Code:         param.Code,
	})
	if err != nil {
		Err := errors.ErrWriteError.Wrap(err, "unable to persist the refresh token")
		o.logger.Error(ctx, "error saving the refresh token", zap.Error(Err), zap.Any("refresh-token", param))
		return nil, Err
	}
	return &dto.RefreshToken{
		Code:         refToken.Code,
		RefreshToken: refToken.RefreshToken,
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

func (o *oauth2) RemoveRefreshTokenCode(ctx context.Context, code string) error {
	if err := o.db.RemoveRefreshTokenByCode(ctx, code); err != nil {
		err := errors.ErrDBDelError.Wrap(err, "could be able to delete the referesh token")
		o.logger.Error(ctx, "unable to delete the refresh token", zap.Error(err), zap.Any("refresh-token-code", code))
		return err
	}
	return nil
}
func (o *oauth2) RemoveRefreshToken(ctx context.Context, refresh_token string) error {
	if err := o.db.RemoveRefreshToken(ctx, refresh_token); err != nil {
		err := errors.ErrDBDelError.Wrap(err, "could be able to delete the referesh token")
		o.logger.Error(ctx, "unable to delete the refresh token", zap.Error(err), zap.Any("refresh-token", refresh_token))
		return err
	}
	return nil
}

func (o *oauth2) CheckIfUserGrantedClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) (bool, dto.RefreshToken, error) {
	refereshToken, err := o.db.GetRefreshTokenByUserIDAndClientID(ctx, db.GetRefreshTokenByUserIDAndClientIDParams{
		UserID:   userID,
		ClientID: clientID,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, dto.RefreshToken{}, nil
		}
		err = errors.ErrReadError.Wrap(err, "could not read refresh token")
		o.logger.Error(ctx, "error could not check if user granted", zap.Error(err), zap.Any("user-id", userID), zap.Any("client-id", clientID))
		return false, dto.RefreshToken{}, err
	}

	return true, dto.RefreshToken{
		ID:           refereshToken.ID,
		Code:         refereshToken.Code,
		RefreshToken: refereshToken.RefreshToken,
		RedirectUri:  refereshToken.RedirectUri.String,
		Scope:        refereshToken.Scope.String,
		UserID:       refereshToken.UserID,
		ClientID:     refereshToken.ClientID,
	}, nil
}

func (o *oauth2) GetRefreshToken(ctx context.Context, token string) (*dto.RefreshToken, error) {
	refreshToken, err := o.db.GetRefreshToken(ctx, token)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no refresh token found")
			o.logger.Info(ctx, "refresh token not found", zap.Error(err), zap.Any("refresh-token", token))
			return nil, err
		}
		err = errors.ErrReadError.Wrap(err, "could not read refresh token")
		o.logger.Error(ctx, "could not found refresh token", zap.Error(err))
		return nil, err
	}
	return &dto.RefreshToken{
		ID:           refreshToken.ID,
		Code:         refreshToken.Code,
		RefreshToken: refreshToken.RefreshToken,
		RedirectUri:  refreshToken.RedirectUri.String,
		Scope:        refreshToken.Scope.String,
		UserID:       refreshToken.UserID,
		ClientID:     refreshToken.ClientID,
		ExpiresAt:    refreshToken.ExpiresAt,
	}, nil
}

func (o *oauth2) GetRefreshTokenOfClientByUserID(ctx context.Context, userID, clientID uuid.UUID) (*dto.RefreshToken, error) {
	refreshToken, err := o.db.GetRefreshTokenByUserIDAndClientID(ctx, db.GetRefreshTokenByUserIDAndClientIDParams{
		UserID:   userID,
		ClientID: clientID,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no refresh token found")
			o.logger.Info(ctx, "refresh token not found", zap.Error(err), zap.Any("user-id", userID), zap.Any("client-id", clientID))
			return nil, err
		}
		err = errors.ErrReadError.Wrap(err, "could not read refresh token")
		o.logger.Error(ctx, "could not find refresh token", zap.Error(err), zap.Any("user-id", userID), zap.Any("client-id", clientID))
		return nil, err
	}
	return &dto.RefreshToken{
		ID:           refreshToken.ID,
		Code:         refreshToken.Code,
		RefreshToken: refreshToken.RefreshToken,
		RedirectUri:  refreshToken.RedirectUri.String,
		Scope:        refreshToken.Scope.String,
		UserID:       refreshToken.UserID,
		ClientID:     refreshToken.ClientID,
		ExpiresAt:    refreshToken.ExpiresAt,
	}, nil
}

func (o *oauth2) GetAuthorizedClients(ctx context.Context, userID uuid.UUID) ([]dto.AuthorizedClientsResponse, error) {
	authorizedClients, err := o.db.GetAuthorizedClientsForUser(ctx, userID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no authorized clients found")
			o.logger.Info(ctx, "no authorized clients were found", zap.Error(err), zap.Any("user-id", userID))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading authorized clients")
			o.logger.Error(ctx, "error reading authorized clients", zap.Error(err), zap.Any("user-id", userID))
			return nil, err
		}
	}
	authorizedClientsDTO := make([]dto.AuthorizedClientsResponse, len(authorizedClients))
	for k, v := range authorizedClients {
		var scopes []dto.Scope
		for _, s := range utils.StringToArray(v.Scope.String) {
			if s == "openid" {
				continue
			}
			scope, err := o.db.GetScope(ctx, s)
			if err != nil {
				if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
					err := errors.ErrNoRecordFound.Wrap(err, "scope doesn't exist")
					o.logger.Error(ctx, "scope row for given scope was not found", zap.Error(err), zap.Any("user-id", userID), zap.String("scope", s))
					continue
				} else {
					err = errors.ErrReadError.Wrap(err, "error reading scope")
					o.logger.Error(ctx, "error encountered while reading scope for authorized client", zap.Error(err), zap.Any("user-id", userID), zap.String("scope", s))
					return nil, err
				}
			}
			scopes = append(scopes, dto.Scope{
				Name:               scope.Name,
				Description:        scope.Description,
				ResourceServerName: scope.ResourceServerName.String,
			})
		}
		authorizedClientsDTO[k] = dto.AuthorizedClientsResponse{
			Client: dto.Client{
				ID:         v.ID,
				Name:       v.Name,
				ClientType: v.ClientType,
				LogoURL:    v.LogoUrl,
			},
			AuthGivenAt:   v.CreatedAt,
			AuthUpdatedAt: v.UpdatedAt,
			AuthExpiresAt: v.ExpiresAt,
			AuthScopes:    scopes,
		}
	}
	return authorizedClientsDTO, nil
}

func (o *oauth2) GetOpenIDAuthorizedClients(ctx context.Context, userID uuid.UUID) ([]dto.AuthorizedClientsResponse, error) {
	authorizedClients, err := o.db.GetOpenIDAuthorizedClientsForUser(ctx, userID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no openid authorized clients found")
			o.logger.Info(ctx, "no openid authorized clients were found", zap.Error(err), zap.Any("user-id", userID))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading openid authorized clients")
			o.logger.Error(ctx, "error reading openid authorized clients", zap.Error(err), zap.Any("user-id", userID))
			return nil, err
		}
	}
	authorizedClientsDTO := make([]dto.AuthorizedClientsResponse, len(authorizedClients))
	for k, v := range authorizedClients {
		authorizedClientsDTO[k] = dto.AuthorizedClientsResponse{
			Client: dto.Client{
				ID:         v.ID,
				Name:       v.Name,
				ClientType: v.ClientType,
				LogoURL:    v.LogoUrl,
			},
			AuthGivenAt:   v.CreatedAt,
			AuthUpdatedAt: v.UpdatedAt,
			AuthExpiresAt: v.ExpiresAt,
		}
	}
	return authorizedClientsDTO, nil
}

func (o *oauth2) UserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	user, err := o.db.GetUserById(ctx, userID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user were found", zap.Error(err), zap.Any("user-id", userID))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading users")
			o.logger.Error(ctx, "error reading users", zap.Error(err), zap.Any("user-id", userID))
			return nil, err
		}
	}

	return &dto.UserInfo{
		Sub:        user.ID.String(),
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Gender:     user.Gender,
		Email:      user.Email.String,
		Phone:      user.Phone,
	}, nil
}

func (o *oauth2) UpdateRefreshToken(ctx context.Context, newRefreshToken, oldRefreshToken string) (*dto.RefreshToken, error) {
	refreshToken, err := o.db.UpdateOAuthRefreshToken(ctx, db.UpdateOAuthRefreshTokenParams{
		RefreshToken:   newRefreshToken,
		RefreshToken_2: oldRefreshToken,
	})

	if err != nil {
		err := errors.ErrUpdateError.Wrap(err, "error updating refresh token")
		o.logger.Error(ctx, "error while updating refresh token for a client access token grant",
			zap.Any("old-refresh-token", oldRefreshToken),
			zap.Any("new-refresh-token", newRefreshToken))
		return nil, err
	}

	return &dto.RefreshToken{
		ID:           refreshToken.ID,
		RefreshToken: refreshToken.RefreshToken,
		Code:         refreshToken.Code,
		UserID:       refreshToken.UserID,
		ClientID:     refreshToken.ClientID,
		Scope:        refreshToken.Scope.String,
		RedirectUri:  refreshToken.RedirectUri.String,
		ExpiresAt:    refreshToken.ExpiresAt,
		CreatedAt:    refreshToken.CreatedAt,
	}, nil
}
