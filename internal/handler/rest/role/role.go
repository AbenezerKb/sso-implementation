package role

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
)

type role struct {
	logger     logger.Logger
	roleModule module.RoleModule
}

func InitRole(logger logger.Logger, roleModule module.RoleModule) rest.Role {
	return &role{
		logger:     logger,
		roleModule: roleModule,
	}
}

// GetAllPermissions is used to get the list of predefined permissions
// @Summary Get all permissions
// @Description Get all permissions that are predefined and fixed for this server
// @ID get-all-permissions
// @Tags role
// @Accept  json
// @Produce  json
// @Param category query string true "category of permissions"
// @Success 200 {object} permissions.Permission
// @Failure 400 {object} model.ErrorResponse
// @Router /roles/permissions [get]
// @Security BearerAuth
func (r *role) GetAllPermissions(ctx *gin.Context) {
	var category request_models.PermissionCategory
	err := ctx.BindQuery(&category)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "could not bind to PermissionCategory. invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	perms, err := r.roleModule.GetAllPermissions(ctx, category.Category)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, perms, nil)
}

// CreateRole is used to create a role with specified permissions
// @Summary create role
// @Description create a role with specified name and permission list
// @ID create-role
// @Tags role
// @Accept  json
// @Produce  json
// @Param role body dto.Role true "role"
// @Success 200 {object} dto.Role
// @Failure 400 {object} model.ErrorResponse
// @Router /roles [post]
// @Security BearerAuth
func (r *role) CreateRole(ctx *gin.Context) {
	var role dto.Role
	err := ctx.ShouldBind(&role)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "could not bind to dto.Role. invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	roleCreated, err := r.roleModule.CreateRole(ctx, role)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusCreated, roleCreated, nil)
}

// GetAllRoles returns all roles
// @Summary      returns all roles that satisfy the given filters
// @Description  returns all roles based on the filters and pagination given
// @Tags         role
// @Accept       json
// @Produce      json
// @param filter query request_models.PgnFltQueryParams true "filter"
// @Success      200  {object}  []dto.Role
// @Failure      400  {object}  model.ErrorResponse
// @Router       /roles [get]
// @Security	BearerAuth
func (r *role) GetAllRoles(ctx *gin.Context) {
	var filtersParam request_models.PgnFltQueryParams
	err := ctx.BindQuery(&filtersParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		r.logger.Info(ctx, "invalid query params", zap.Error(err), zap.Any("query-params", ctx.Request.URL.Query()))
		_ = ctx.Error(err)
		return
	}

	roles, metaData, err := r.roleModule.GetAllRoles(ctx, filtersParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, roles, metaData)
}

// UpdateRoleStatus updates role status
// @Summary      changes role status
// @Description  changes role status so that they can revoke a role
// @Tags         role
// @Accept       json
// @Produce      json
// @param status body dto.UpdateRoleStatus true "status"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.ErrorResponse
// @Router       /roles/{name}/status [patch]
// @Security	BearerAuth
func (r *role) UpdateRoleStatus(ctx *gin.Context) {
	roleName := ctx.Param("name")
	updateStatusParam := dto.UpdateRoleStatus{}
	err := ctx.ShouldBindJSON(&updateStatusParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "unable to bind role status", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	err = r.roleModule.UpdateRoleStatus(requestCtx, updateStatusParam, roleName)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	r.logger.Info(ctx, "role status changed", zap.String("role-name", roleName), zap.String("to-status", updateStatusParam.Status))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// DeleteRole deletes a role
// @Summary      deletes a role
// @Description  deletes a role and all user associations with the role
// @Tags         role
// @Accept       json
// @Produce      json
// @param name path string true "name"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.ErrorResponse
// @Router       /roles/{name} [delete]
// @Security	BearerAuth
func (r *role) DeleteRole(ctx *gin.Context) {
	roleName := ctx.Param("name")
	requestCtx := ctx.Request.Context()
	err := r.roleModule.DeleteRole(requestCtx, roleName)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	r.logger.Info(ctx, "role deleted", zap.String("role-name", roleName))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// UpdateRole updates a role
// @Summary      updates a role
// @Description  updates a role with new permissions
// @Tags         role
// @Accept       json
// @Produce      json
// @param role body dto.UpdateRole true "body"
// @Success      200  {object}  dto.Role
// @Failure      400  {object}  model.ErrorResponse
// @Router       /roles/{name} [put]
// @Security	BearerAuth
func (r *role) UpdateRole(ctx *gin.Context) {
	roleName := ctx.Param("name")
	requestCtx := ctx.Request.Context()
	var updateRole dto.UpdateRole
	err := ctx.ShouldBind(&updateRole)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.logger.Info(ctx, "invalid input on update role", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	updateRole.Name = roleName
	role, err := r.roleModule.UpdateRole(requestCtx, updateRole)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	r.logger.Info(ctx, "role updated", zap.String("role-name", roleName))
	constant.SuccessResponse(ctx, http.StatusOK, role, nil)
}
