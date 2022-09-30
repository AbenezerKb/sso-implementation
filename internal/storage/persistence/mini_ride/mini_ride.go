package mini_ride

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"go.uber.org/zap"
)

type miniRidePersistence struct {
	logger logger.Logger
	db     *persistencedb.PersistenceDB
}

func InitMiniRidePersistence(logger logger.Logger, db *persistencedb.PersistenceDB) storage.MiniRidePersistence {
	return &miniRidePersistence{
		logger: logger,
		db:     db,
	}
}

func (u *miniRidePersistence) UpdateUser(ctx context.Context, updateUserParam *request_models.Driver) error {
	_, err := u.db.UpdateUserByID(ctx, db.UpdateUserByIDParams{
		FirstName:      updateUserParam.FirstName,
		MiddleName:     updateUserParam.MiddleName,
		LastName:       updateUserParam.LastName,
		Status:         sql.NullString{String: updateUserParam.Status, Valid: true},
		ProfilePicture: sql.NullString{String: updateUserParam.ProfilePicture, Valid: true},
		Phone:          updateUserParam.Phone,
		ID:             updateUserParam.ID,
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not update user profile")
		u.logger.Error(ctx, "unable to update user profile", zap.Error(err), zap.Any("user", updateUserParam))
		return err
	}

	return nil
}

func (m *miniRidePersistence) CreateUser(ctx context.Context, createUserParam *request_models.Driver) (*dto.User, error) {
	registeredUser, err := m.db.CreateUserWithID(ctx, db.CreateUserWithIDParams{
		FirstName:      createUserParam.FirstName,
		LastName:       createUserParam.LastName,
		Gender:         createUserParam.Gender,
		MiddleName:     createUserParam.MiddleName,
		ProfilePicture: utils.StringOrNull(createUserParam.ProfilePicture),
		Phone:          createUserParam.Phone,
		ID:             createUserParam.ID,
	})
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create user")
		m.logger.Error(ctx, "unable to create user from mini-ride", zap.Error(err), zap.Any("user", createUserParam))
		return nil, err
	}
	return &dto.User{
		ID:             registeredUser.ID,
		Status:         registeredUser.Status.String,
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

func (u *miniRidePersistence) SwapPhones(ctx context.Context, newPhone, oldPhone string) error {
	err := u.db.SwapPhones(ctx, newPhone, oldPhone)
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "error swapping phone")
		u.logger.Error(ctx, "couldn't swap phone", zap.Error(err), zap.String("phone1", newPhone), zap.String("phone2", oldPhone))
		return err
	}
	return nil
}
