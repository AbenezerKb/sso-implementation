package role

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
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
// @Param category body string true "category of permissions"
// @Success 200 {object} permissions.Permission
// @Failure 400 {object} model.ErrorResponse
// @Router /roles/permissions [get]
// @Security BearerAuth
func (r *role) GetAllPermissions(ctx *gin.Context) {
	var category request_models.PermissionCategory
	err := ctx.ShouldBind(&category)
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
