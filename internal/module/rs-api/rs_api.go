package rs_api

import (
	"context"

	"sso/internal/constant/errors"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"

	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
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
	filtersParam db_pgnflt.PgnFltQueryParams) (*dto.RSAPIUsersResponse, *model.MetaData, error) {
	if err := req.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("request", req))
		return nil, nil, err
	}

	filters, err := filtersParam.ToFilterParams([]db_pgnflt.FieldType{
		{Name: "first_name", Type: db_pgnflt.String},
		{Name: "middle_name", Type: db_pgnflt.String},
		{Name: "last_name", Type: db_pgnflt.String},
		{Name: "email", Type: db_pgnflt.String},
		{Name: "phone", Type: db_pgnflt.String},
		{Name: "gender", Type: db_pgnflt.String},
		{Name: "status", Type: db_pgnflt.Enum,
			Values: []string{"ACTIVE", "INACTIVE", "PENDING"},
		},
		{Name: "created_at", Type: db_pgnflt.Time},
	}, db_pgnflt.Defaults{
		Sort: []db_pgnflt.Sort{
			{
				Field: "created_at",
				Sort:  db_pgnflt.SortDesc,
			},
		},
		PerPage: 10,
	})
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid filter params")
		r.logger.Info(ctx, "invalid filter params were given", zap.Error(err), zap.Any("filters-query", filtersParam))
		return nil, nil, err
	}

	var res dto.RSAPIUsersResponse
	var total int
	// get users by phones and ids
	if len(req.IDs) > 0 {
		// fetch users by id
		usersPart, metaData, err := r.userPersistence.GetUsersByID(ctx, req.IDs, filters)
		if err != nil {
			return nil, nil, err
		}
		res.IDs = usersPart
		total += metaData.Total
	}

	if len(req.Phones) > 0 {
		// fetch users by phone
		var parsedPhones []string
		for i := 0; i < len(req.Phones); i++ {
			parsedPhones = append(parsedPhones, phonenumber.Parse(req.Phones[i], "ET"))
		}
		usersPart, metaData, err := r.userPersistence.GetUsersByPhone(ctx, parsedPhones, filters)
		if err != nil {
			return nil, nil, err
		}
		res.Phones = usersPart
		total += metaData.Total
	}

	return &res, &model.MetaData{
		FilterParams: filters,
		Total:        total,
	}, nil
}
