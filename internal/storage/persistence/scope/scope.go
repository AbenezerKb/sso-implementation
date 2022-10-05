package scope

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"

	"go.uber.org/zap"
)

type scopePersistence struct {
	logger logger.Logger
	db     *db.Queries
}

func InitScopePersistence(logger logger.Logger, db *db.Queries) storage.ScopePersistence {
	return &scopePersistence{
		logger: logger,
		db:     db,
	}
}

func (s *scopePersistence) GetListedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error) {
	listedScopes := []dto.Scope{}
	for _, scope := range scopes {
		scope, err := s.db.GetScope(ctx, scope)
		if err != nil {
			if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
				continue
			}
			err = errors.ErrReadError.Wrap(err, "could not read the scope")
			s.logger.Error(ctx, "unable to read the scope", zap.Error(err), zap.Any("scope", scope))
			return nil, err
		}
		listedScopes = append(listedScopes, dto.Scope{
			Name:        scope.Name,
			Description: scope.Description,
		})
	}
	return listedScopes, nil
}

func (s *scopePersistence) GetScope(ctx context.Context, scope string) (dto.Scope, error) {
	createdScope, err := s.db.GetScope(ctx, scope)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return dto.Scope{}, errors.ErrNoRecordFound.Wrap(err, "scope not found")
		}
		err = errors.ErrReadError.Wrap(err, "could not read the scope")
		s.logger.Error(ctx, "unable to read the scope", zap.Error(err), zap.Any("scope", scope))
		return dto.Scope{}, err
	}
	return dto.Scope{
		Name:        createdScope.Name,
		Description: createdScope.Description,
	}, nil
}

func (s *scopePersistence) CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error) {
	createdScope, err := s.db.CreateScope(ctx, db.CreateScopeParams{
		Name:               scope.Name,
		Description:        scope.Description,
		ResourceServerName: sql.NullString{String: scope.ResourceServerName, Valid: true},
	})

	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create the scope")
		s.logger.Error(ctx, "unable to create the scope", zap.Error(err), zap.Any("scope", scope))
		return dto.Scope{}, err
	}
	return dto.Scope{
		Name:        createdScope.Name,
		Description: createdScope.Description,
	}, nil
}

func (s *scopePersistence) GetScopeNameOnly(ctx context.Context, scopes ...string) (string, error) {
	scopesAry, err := s.GetListedScopes(ctx, scopes...)
	if err != nil {
		return "", err
	}
	scopeNameAry := []string{}
	for _, x := range scopesAry {
		scopeNameAry = append(scopeNameAry, x.Name)
	}
	scopeStr := utils.ArrayToString(scopeNameAry)
	return scopeStr, nil
}

func (s *scopePersistence) GetAllScopes(ctx context.Context, filters request_models.FilterParams) ([]dto.Scope, *model.MetaData, error) {
	scopes, total, err := s.db.GetAllScopes(ctx, utils.ComposeFilterSQL(ctx, filters, s.logger))
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no scope found")
			s.logger.Info(ctx, "no scopes were found", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading scopes")
			s.logger.Error(ctx, "error reading scopes", zap.Error(err), zap.Any("filters", filters))
			return nil, nil, err
		}
	}
	scopesDTO := make([]dto.Scope, len(scopes))
	for k, v := range scopes {
		scopesDTO[k] = dto.Scope{
			Description:        v.Description,
			Name:               v.Name,
			CreatedAt:          v.CreatedAt,
			ResourceServerName: v.ResourceServerName.String,
		}
	}
	return scopesDTO, &model.MetaData{
		FilterParams: filters,
		Total:        total,
		Extra:        nil,
	}, nil
}

func (s *scopePersistence) DeleteScopeByName(ctx context.Context, name string) error {
	_, err := s.db.DeleteScope(ctx, name)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no scope found")
			s.logger.Info(ctx, "no scopes were found", zap.Error(err), zap.String("scope", name))
			return err
		} else {
			err = errors.ErrDBDelError.Wrap(err, "error deleting scope")
			s.logger.Error(ctx, "error deleting scope", zap.Error(err), zap.String("scope", name))
			return err
		}
	}

	return nil
}
