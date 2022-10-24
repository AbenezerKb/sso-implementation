package oauth

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

type oauth struct {
	logger logger.Logger
	db     *db.Queries
}

func InitOAuth(logger logger.Logger, db *db.Queries) storage.OAuthPersistence {
	return &oauth{
		logger,
		db,
	}
}

func (o *oauth) Register(ctx context.Context, userParam dto.User) (*dto.User, error) {
	registeredUser, err := o.db.CreateUser(ctx, db.CreateUserParams{
		FirstName:      userParam.FirstName,
		LastName:       userParam.LastName,
		Email:          utils.StringOrNull(userParam.Email),
		Gender:         userParam.Gender,
		MiddleName:     userParam.MiddleName,
		ProfilePicture: utils.StringOrNull(userParam.ProfilePicture),
		Phone:          userParam.Phone,
		Password:       userParam.Password,
	})
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create user")
		o.logger.Error(ctx, "unable to create user", zap.Error(err), zap.Any("user", userParam))
		return nil, err
	}
	return &dto.User{
		ID:             registeredUser.ID,
		Status:         registeredUser.Status.String,
		UserName:       registeredUser.UserName,
		FirstName:      registeredUser.FirstName,
		MiddleName:     registeredUser.MiddleName,
		LastName:       registeredUser.LastName,
		Email:          registeredUser.Email.String,
		Phone:          registeredUser.Phone,
		Gender:         registeredUser.Gender,
		CreatedAt:      registeredUser.CreatedAt,
		ProfilePicture: registeredUser.ProfilePicture.String,
	}, nil
}

func (o *oauth) GetUserByPhone(ctx context.Context, phone string) (*dto.User, error) {
	user, err := o.db.GetUserByPhone(ctx, phone)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user found", zap.Error(err), zap.String("phone", phone))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by phone", zap.Error(err), zap.String("phone", phone))
			return nil, err
		}
	}
	return &dto.User{
		ID:             user.ID,
		Status:         user.Status.String,
		UserName:       user.UserName,
		FirstName:      user.FirstName,
		MiddleName:     user.MiddleName,
		LastName:       user.LastName,
		Email:          user.Email.String,
		Phone:          user.Phone,
		Gender:         user.Gender,
		ProfilePicture: user.ProfilePicture.String,
		Password:       user.Password,
	}, nil
}
func (o *oauth) GetUserByEmail(ctx context.Context, email string) (*dto.User, error) {
	user, err := o.db.GetUserByEmail(ctx, utils.StringOrNull(email))
	if err != nil {

		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user found by email", zap.Error(err), zap.String("email", email))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by email", zap.Error(err), zap.String("email", email))
			return nil, err
		}
	}

	return &dto.User{
		ID:             user.ID,
		Status:         user.Status.String,
		UserName:       user.UserName,
		FirstName:      user.FirstName,
		MiddleName:     user.MiddleName,
		LastName:       user.LastName,
		Email:          user.Email.String,
		Phone:          user.Phone,
		Gender:         user.Gender,
		ProfilePicture: user.ProfilePicture.String,
	}, nil
}

func (o *oauth) GetUserStatus(ctx context.Context, Id uuid.UUID) (string, error) {

	status, err := o.db.GetUserStatus(ctx, Id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user found to fetch user status", zap.Error(err), zap.String("id", Id.String()))
			return "", err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by id", zap.Error(err), zap.String("id", Id.String()))
			return "", err
		}
	}

	return status.String, nil
}
func (o *oauth) UserByPhoneExists(ctx context.Context, phone string) (bool, error) {
	_, err := o.db.GetUserByPhone(ctx, phone)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by phone", zap.Error(err), zap.String("phone", phone))
			return false, err
		}
	}
	return true, nil
}

func (o *oauth) GetUserByPhoneOrEmail(ctx context.Context, query string) (*dto.User, error) {
	user, err := o.db.GetUserByPhoneOrEmail(ctx, query)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user found", zap.Error(err), zap.String("query", query))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by phone or email", zap.Error(err), zap.String("email-or-phone", query))
			return nil, err
		}
	}

	return &dto.User{
		ID:         user.ID,
		Status:     user.Status.String,
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email:      user.Email.String,
		Phone:      user.Phone,
		Password:   user.Password,
	}, nil
}

func (o *oauth) UserByEmailExists(ctx context.Context, email string) (bool, error) {
	_, err := o.db.GetUserByEmail(ctx, utils.StringOrNull(email))
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by email", zap.Error(err), zap.String("email", email))
			return false, err
		}
	}
	return true, nil
}

func (o *oauth) GetUserByID(ctx context.Context, Id uuid.UUID) (*dto.User, error) {
	user, err := o.db.GetUserById(ctx, Id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user found", zap.Error(err), zap.String("id", Id.String()))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by id", zap.Error(err), zap.String("id", Id.String()))
			return nil, err
		}
	}

	return &dto.User{
		ID:         user.ID,
		Status:     user.Status.String,
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email:      user.Email.String,
		Phone:      user.Phone,
	}, nil
}

func (o *oauth) SaveInternalRefreshToken(ctx context.Context, rf dto.InternalRefreshToken) error {
	_, err := o.db.SaveInternalRefreshToken(ctx, db.SaveInternalRefreshTokenParams{
		UserID:       rf.UserID,
		RefreshToken: rf.RefreshToken,
		IpAddress:    rf.IPAddress,
		UserAgent:    rf.UserAgent,
		ExpiresAt:    rf.ExpiresAt,
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not save internal rf token")
		o.logger.Error(ctx, "could not save internal refresh token", zap.Error(err), zap.Any("internalRefrshToken", rf))
		return err
	}

	return nil
}

func (o *oauth) RemoveInternalRefreshToken(ctx context.Context, refreshToken string) error {
	err := o.db.RemoveInternalRefreshToken(ctx, refreshToken)
	if err != nil {
		err = errors.ErrDBDelError.Wrap(err, "could not remove internal rf token")
		o.logger.Error(ctx, "could not remove internal rf token", zap.Error(err), zap.Any("refresh-token", refreshToken))
		return err
	}

	return nil
}

func (o *oauth) GetInternalRefreshToken(ctx context.Context, token string) (*dto.InternalRefreshToken, error) {
	refreshToken, err := o.db.GetInternalRefreshToken(ctx, token)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no refresh token found")
			o.logger.Info(ctx, "internal refresh token not found", zap.Error(err), zap.Any("internal-refresh-token", token))
			return nil, err
		}
		err = errors.ErrReadError.Wrap(err, "could not read refresh token")
		o.logger.Error(ctx, "could not found refresh token", zap.Error(err))
		return nil, err
	}
	return &dto.InternalRefreshToken{
		ID:           refreshToken.ID,
		RefreshToken: refreshToken.RefreshToken,
		ExpiresAt:    refreshToken.ExpiresAt,
		UserAgent:    refreshToken.UserAgent,
		IPAddress:    refreshToken.IpAddress,
		UserID:       refreshToken.UserID,
		CreatedAt:    refreshToken.CreatedAt,
	}, nil
}

func (o *oauth) UpdateInternalRefreshToken(ctx context.Context, param dto.InternalRefreshToken) (*dto.InternalRefreshToken, error) {
	refreshToken, err := o.db.UpdateRefreshToken(ctx, db.UpdateRefreshTokenParams{
		ExpiresAt:    param.ExpiresAt,
		RefreshToken: param.RefreshToken,
		ID:           param.ID,
	})
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "unable to update the refresh token")
		o.logger.Error(ctx, "error updating the user refresh ytoken", zap.Error(err), zap.String("internal-refresh-token", param.ID.String()))
		return nil, err
	}
	return &dto.InternalRefreshToken{
		ID:           refreshToken.ID,
		RefreshToken: refreshToken.RefreshToken,
		ExpiresAt:    refreshToken.ExpiresAt,
		UserID:       refreshToken.UserID,
		CreatedAt:    refreshToken.CreatedAt,
	}, nil
}

func (o *oauth) GetInternalRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]dto.InternalRefreshToken, error) {
	refreshTokens, err := o.db.GetInternalRefreshTokensByUserID(ctx, userID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no refresh token found")
			o.logger.Info(ctx, "internal refresh token not found", zap.Error(err), zap.Any("internal-refresh-token", userID))
			return nil, err
		}
		err = errors.ErrReadError.Wrap(err, "could not read refresh token")
		o.logger.Error(ctx, "could not found refresh token", zap.Error(err))
		return nil, err
	}

	dtoRefreshTokens := make([]dto.InternalRefreshToken, len(refreshTokens))
	for i := 0; i < len(refreshTokens); i++ {
		dtoRefreshTokens[i] = dto.InternalRefreshToken{
			ID:           refreshTokens[i].ID,
			RefreshToken: refreshTokens[i].RefreshToken,
			ExpiresAt:    refreshTokens[i].ExpiresAt,
			UserID:       refreshTokens[i].UserID,
			CreatedAt:    refreshTokens[i].CreatedAt,
		}
	}

	return dtoRefreshTokens, nil
}

func (o *oauth) GetUserPassword(ctx context.Context, Id uuid.UUID) (string, error) {
	user, err := o.db.GetUserById(ctx, Id)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err = errors.ErrNoRecordFound.Wrap(err, "no user found")
			o.logger.Info(ctx, "no user found", zap.Error(err), zap.String("id", Id.String()))
			return "", err
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, "unable to get user by id", zap.Error(err), zap.String("id", Id.String()))
			return "", err
		}
	}

	return user.Password, nil
}
