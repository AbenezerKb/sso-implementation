package scope

import (
	"context"
	"database/sql"
	"sso/internal/constant/errors"
	"sso/internal/constant/errors/sqlcerr"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/storage"
	"sso/platform/logger"

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
