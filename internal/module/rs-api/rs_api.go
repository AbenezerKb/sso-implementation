package rs_api

import (
	"context"
	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
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
