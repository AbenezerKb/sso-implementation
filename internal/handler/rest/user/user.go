package user

import (
	"fmt"
	"net/http"

	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type user struct {
	logger     logger.Logger
	userModule module.UserModule
}

func Init(logger logger.Logger, userModule module.UserModule) rest.User {
	return &user{
		logger,
		userModule,
	}
}

// CreateUser	 creates a new user.
// @Summary      create a new user.
// @Description  create a new user.
// @Tags         user
// @Accept       json
// @Produce      json
// @param user body dto.CreateUser true "user"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users [post]
// @Security	BearerAuth
func (u *user) CreateUser(ctx *gin.Context) {
	userParam := dto.CreateUser{}
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		u.logger.Info(ctx, "unable to bind user data", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	createdUser, err := u.userModule.Create(ctx.Request.Context(), userParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	u.logger.Info(ctx, "created user")
	constant.SuccessResponse(ctx, http.StatusCreated, createdUser, nil)
}

// GetUser	 get user details.
// @Summary      get user details.
// @Description  get user details.
// @Tags         user
// @Accept       json
// @Produce      json
// @param id path string true "id"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /users/{id} [get]
// @Security	BearerAuth
func (u *user) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	requestCtx := ctx.Request.Context()
	user, err := u.userModule.GetUserByID(requestCtx, userID)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	u.logger.Info(ctx, "user details fetched", zap.Any("user-id", userID))
	constant.SuccessResponse(ctx, http.StatusOK, user, nil)
}

// GetAllUsers returns all users
// @Summary      returns all users that satisfy the given filters
// @Description  returns all users based on the filters and pagination given
// @Tags         user
// @Accept       json
// @Produce      json
// @param filter query request_models.PgnFltQueryParams true "filter"
// @Success      200  {object}  []dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users [get]
// @Security	BearerAuth
func (u *user) GetAllUsers(ctx *gin.Context) {
	var filtersParam request_models.PgnFltQueryParams
	err := ctx.BindQuery(&filtersParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		u.logger.Info(ctx, "invalid query params", zap.Error(err), zap.Any("query-params", ctx.Request.URL.Query()))
		_ = ctx.Error(err)
		return
	}

	users, metaData, err := u.userModule.GetAllUsers(ctx, filtersParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, users, metaData)
}

// UpdateUserStatus updates user status
// @Summary      changes user status
// @Description  changes user status to ACTIVE or INACTIVE
// @Tags         user
// @Accept       json
// @Produce      json
// @param user body dto.CreateUser true "user"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users/{id}/status [patch]
// @Security	BearerAuth
func (u *user) UpdateUserStatus(ctx *gin.Context) {
	userID := ctx.Param("id")
	updateUserStatusParam := dto.UpdateUserStatus{}
	err := ctx.ShouldBindJSON(&updateUserStatusParam)
	fmt.Println(ctx.Request.Body)

	if err != nil {
		u.logger.Info(ctx, "unable to bind user data", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	requestCtx := ctx.Request.Context()

	err = u.userModule.UpdateUserStatus(requestCtx, updateUserStatusParam, userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	u.logger.Info(ctx, "user status changed", zap.Any("user-id", userID), zap.Any("to-status", updateUserStatusParam))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// UpdateUserRole	 updates the role for the user
// @Summary      updates the role for the user
// @Description  updates the role for the user
// @Tags         user
// @Accept       json
// @Produce      json
// @param role body dto.AssignRole true "role"
// @Success      200
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users/{id}/role [patch]
// @Security	BearerAuth
func (u *user) UpdateUserRole(ctx *gin.Context) {
	userID := ctx.Param("id")
	role := dto.AssignRole{}
	err := ctx.ShouldBind(&role)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.logger.Info(ctx, "unable to bind to AssignRole for update user role", zap.Error(err), zap.String("user-id", userID))
		_ = ctx.Error(err)
		return
	}
	err = u.userModule.UpdateUserRole(ctx.Request.Context(), userID, role)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	u.logger.Info(ctx, "updated role for user", zap.String("user-id", userID))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// RevokeUserRole	 revokes the role from the user
// @Summary      revokes the role from the user
// @Description  revokes the role from the user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users/{id}/role [delete]
// @Security	BearerAuth
func (u *user) RevokeUserRole(ctx *gin.Context) {
	userID := ctx.Param("id")
	err := u.userModule.RevokeUserRole(ctx.Request.Context(), userID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	u.logger.Info(ctx, "revoked role from user", zap.String("user-id", userID))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// ResetUserPassword	 revokes the role from the user
// @Summary      resets user password
// @Description  resets user password
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users/{id}/password [patch]
// @Security	BearerAuth
func (u *user) ResetUserPassword(ctx *gin.Context) {
	userID := ctx.Param("id")

	err := u.userModule.ResetUserPassword(ctx.Request.Context(), userID)
	if err != nil {
		_ = ctx.Error(err)

		return
	}

	u.logger.Info(ctx, "user password was reset by admin", zap.String("user-id", userID))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
