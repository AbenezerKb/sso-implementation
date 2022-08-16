package oauth

import (
	"context"
	"database/sql"
	"reflect"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"

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
		Email:          sql.NullString{String: userParam.Email, Valid: true},
		Gender:         userParam.Gender,
		MiddleName:     userParam.MiddleName,
		ProfilePicture: sql.NullString{String: userParam.ProfilePicture, Valid: true},
		Phone:          userParam.Phone,
		Password:       userParam.Password,
	})
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create user")
		o.logger.Error(ctx, zap.Error(err).String)
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

		if reflect.ValueOf(user).IsZero() {
			return &dto.User{}, errors.ErrNoRecordFound.Wrap(err, "no user found")
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, zap.Error(err).String)
			return &dto.User{}, err
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
	user, err := o.db.GetUserByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {

		if reflect.ValueOf(user).IsZero() {
			return &dto.User{}, errors.ErrNoRecordFound.Wrap(err, "no user found")
		} else {
			err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			o.logger.Error(ctx, zap.Error(err).String)
			return &dto.User{}, err
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
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Error(ctx,"unable to get user status", zap.Error(err))
		return "", err
	}

	return status.String, nil
}
func (o *oauth) UserByPhoneExists(ctx context.Context, phone string) (bool, error) {
	user, err := o.db.GetUserByPhone(ctx, phone)
	if err != nil {
		if reflect.ValueOf(user).IsZero() {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, zap.Error(err).String)
			return false, err
		}
	}
	return true, nil
}

func (o *oauth) GetUserByPhoneOrEmail(ctx context.Context, query string) (*dto.User, error) {
	user, err := o.db.GetUserByPhoneOrEmail(ctx, query)
	if err != nil {
		if reflect.ValueOf(user).IsZero() {
			return &dto.User{}, errors.ErrNoRecordFound.Wrap(err, "no user found")
		} else {
			err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			o.logger.Error(ctx, zap.Error(err).String)
			return &dto.User{}, err
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
	user, err := o.db.GetUserByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		if reflect.ValueOf(user).IsZero() {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, zap.Error(err).String)
			return false, err
		}
	}
	return true, nil
}

func (o *oauth) GetUserByID(ctx context.Context, Id uuid.UUID) (*dto.User, error) {
	user, err := o.db.GetUserById(ctx, Id)
	if err != nil {
		if reflect.ValueOf(user).IsZero() {
			return &dto.User{}, errors.ErrNoRecordFound.Wrap(err, "no user found")
		} else {
			err = errors.ErrReadError.Wrap(err, "could not read user data")
			o.logger.Error(ctx, zap.Error(err).String)
			return &dto.User{}, err
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
