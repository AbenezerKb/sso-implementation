package scope

import (
	"net/http"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"sso/internal/constant"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type scope struct {
	logger      logger.Logger
	scopeModule module.ScopeModule
}

func InitScope(logger logger.Logger, scopeModule module.ScopeModule) rest.Scope {
	return &scope{
		logger:      logger,
		scopeModule: scopeModule,
	}
}

// GetScope is used to get a scope
// @Summary Get a scope
// @Description Get a scope
// @ID get-scope
// @Tags scope
// @Accept  json
// @Produce  json
// @Param name path string true "Scope name"
// @Success 200 {object} dto.Scope
// @Failure 400 {object} model.ErrorResponse
// @Router /oauth/scopes/{name} [get]
// @Security BearerAuth
func (s *scope) GetScope(ctx *gin.Context) {
	scopeName := ctx.Param("name")

	requestCtx := ctx.Request.Context()
	scope, err := s.scopeModule.GetScope(requestCtx, scopeName)
	if err != nil {
		ctx.Error(err)
		return
	}

	s.logger.Info(ctx, "scope found", zap.Any("scope", scope))
	constant.SuccessResponse(ctx, http.StatusOK, scope, nil)
}

// CreateScope is used to create a new scope
// @Summary Create a new scope
// @Description Create a new scope
// @ID create-scope
// @Tags scope
// @Accept  json
// @Produce  json
// @Param scope body dto.Scope true "Create a new scope"
// @Success 201 {object} dto.Scope
// @Failure 400 {object} model.ErrorResponse
// @Router /oauth/scopes [post]
// @Security BearerAuth
func (s *scope) CreateScope(ctx *gin.Context) {
	scopeParam := dto.Scope{}
	err := ctx.ShouldBind(&scopeParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.logger.Info(ctx, "couldn't bind", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	createdScope, err := s.scopeModule.CreateScope(requestCtx, scopeParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	s.logger.Info(ctx, "created scope", zap.Any("scope", createdScope))
	constant.SuccessResponse(ctx, http.StatusCreated, createdScope, nil)
}

// GetAllScopes returns all scopes
// @Summary returns all scopes that satisfy given filters
// @Description returns all scopes that satisfy given filters
// @Tags scope
// @Accept  json
// @Produce  json
// @param filter query request_models.PgnFltQueryParams true "filter"
// @Success 200 {object} []dto.Scope
// @Failure 400 {object} model.ErrorResponse
// @Router /oauth/scopes [get]
// @Security BearerAuth
func (s *scope) GetAllScopes(ctx *gin.Context) {
	var filtersParam request_models.PgnFltQueryParams
	err := ctx.BindQuery(&filtersParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		s.logger.Info(ctx, "invalid query params", zap.Error(err), zap.Any("query-params", ctx.Request.URL.Query()))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()

	scopes, metaData, err := s.scopeModule.GetAllScopes(requestCtx, filtersParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, scopes, metaData)
}
