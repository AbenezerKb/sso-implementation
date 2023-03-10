package rs_api

import (
	"context"

	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type rsAPI struct {
	logger          logger.Logger
	userPersistence storage.UserPersistence
}

func Init(
	logger logger.Logger,
	userPersistence storage.UserPersistence) module.RSAPI {
	return &rsAPI{
		logger:          logger,
		userPersistence: userPersistence,
	}
}

func (r *rsAPI) GetUserByIDOrPhone(ctx context.Context, req request_models.RSAPIUserRequest) (*dto.User, error) {
	if err := req.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("request", req))
		return nil, err
	}

	if req.ID != "" {
		userID, err := uuid.Parse(req.ID)
		if err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			r.logger.Info(ctx, "invalid input", zap.Error(err), zap.String("user-id", req.ID))
			return nil, err
		}

		return r.userPersistence.GetUserByID(ctx, userID)
	} else {
		return r.userPersistence.GetUserByPhone(ctx, phonenumber.Parse(req.Phone, "ET"))
	}
}
func (r *rsAPI) GetUsersByIDOrPhone(ctx context.Context,
	req request_models.RSAPIUsersRequest,
) (*dto.RSAPIUsersResponse, error) {
	if err := req.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("request", req))
		return nil, err
	}

	var res dto.RSAPIUsersResponse

	// get users by phones and ids
	if len(req.IDs) > 0 {
		// fetch users by id
		usersPart, err := r.userPersistence.GetUsersByID(ctx, req.IDs)
		if err != nil {
			return nil, err
		}

		res.IDs = usersPart
	}

	if len(req.Phones) > 0 {
		// fetch users by phone
		var parsedPhones []string
		for i := 0; i < len(req.Phones); i++ {
			parsedPhones = append(parsedPhones, phonenumber.Parse(req.Phones[i], "ET"))
		}

		usersPart, err := r.userPersistence.GetUsersByPhone(ctx, parsedPhones)
		if err != nil {
			return nil, err
		}

		res.Phones = usersPart
	}

	return &res, nil
}
