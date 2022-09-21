package client

import (
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

type client struct {
	logger       logger.Logger
	clientModule module.ClientModule
}

func Init(logger logger.Logger, clientModule module.ClientModule) rest.Client {
	return &client{
		logger:       logger,
		clientModule: clientModule,
	}
}

// CreateClient is a handler for creating a client
// @Summary      Create a client
// @Description  Create a new client
// @Tags         client
// @Accept       json
// @Produce      json
// @param client body dto.Client true "client"
// @Success      200  {object}  dto.Client
// @Failure      400  {object}  model.ErrorResponse
// @Router       /clients [post]
// @Security	BearerAuth
func (c *client) CreateClient(ctx *gin.Context) {
	clientParam := dto.Client{}
	err := ctx.ShouldBind(&clientParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.logger.Info(ctx, "couldn't bind to dto.Client body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	createdClient, err := c.clientModule.Create(ctx.Request.Context(), clientParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	c.logger.Info(ctx, "created client")
	constant.SuccessResponse(ctx, http.StatusCreated, createdClient, nil)
}

// GetAllClients returns all clients
// @Summary      returns all clients that satisfy the given filters
// @Description  returns all clients based on the filters and pagination given
// @Tags         client
// @Accept       json
// @Produce      json
// @param filter query request_models.PgnFltQueryParams true "filter"
// @Success      200  {object}  dto.Client
// @Failure      400  {object}  model.ErrorResponse
// @Router       /clients [get]
// @Security	BearerAuth
func (c *client) GetAllClients(ctx *gin.Context) {
	var filtersParam request_models.PgnFltQueryParams
	err := ctx.BindQuery(&filtersParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		c.logger.Info(ctx, "invalid query params", zap.Error(err), zap.Any("query-params", ctx.Request.URL.Query()))
		_ = ctx.Error(err)
		return
	}

	clients, metaData, err := c.clientModule.GetAllClients(ctx, filtersParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, clients, metaData)
}

// DeleteClient is a handler for deleting a client
// @Summary      Delete  client
// @Description  Delete  client
// @Tags         client
// @Accept       json
// @Produce      json
// @param id path string true "id"
// @Success      204
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /clients/{id} [delete]
// @Security	BearerAuth
func (c *client) DeleteClient(ctx *gin.Context) {
	clientID := ctx.Param("id")

	requestCtx := ctx.Request.Context()
	err := c.clientModule.DeleteClientByID(requestCtx, clientID)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	c.logger.Info(ctx, "client deleted", zap.Any("client-id", clientID))
	constant.SuccessResponse(ctx, http.StatusNoContent, nil, nil)
}
