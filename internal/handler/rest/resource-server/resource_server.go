package resource_server

import (
	"net/http"

	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
	"go.uber.org/zap"
)

type resourceServer struct {
	logger               logger.Logger
	resourceServerModule module.ResourceServerModule
}

func Init(logger logger.Logger, resourceServerModule module.ResourceServerModule) rest.ResourceServer {
	return &resourceServer{
		logger:               logger,
		resourceServerModule: resourceServerModule,
	}
}

// CreateResourceServer is a handler for creating a resource server
// @Summary      Create a resource server
// @Description  Create a new resource server
// @Tags         resourceServer
// @Accept       json
// @Produce      json
// @param client body dto.ResourceServer true "resource_server"
// @Success      200  {object}  dto.ResourceServer
// @Failure      400  {object}  model.ErrorResponse
// @Router       /resourceServers [post]
// @Security	BearerAuth
func (c *resourceServer) CreateResourceServer(ctx *gin.Context) {
	resourceServerBody := dto.ResourceServer{}
	err := ctx.ShouldBind(&resourceServerBody)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.logger.Info(ctx, "couldn't bind to dto.ResourceServer body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	createdServer, err := c.resourceServerModule.CreateResourceServer(ctx.Request.Context(), resourceServerBody)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	c.logger.Info(ctx, "created resource server")
	constant.SuccessResponse(ctx, http.StatusCreated, createdServer, nil)
}

// GetAllResourceServers returns all resource servers
// @Summary      returns all resource servers that satisfy the given filters
// @Description  returns all resource servers based on the filters and pagination given
// @Tags         resourceServer
// @Accept       json
// @Produce      json
// @param filter query request_models.PgnFltQueryParams true "filter"
// @Success      200  {object}  []dto.ResourceServer
// @Failure      400  {object}  model.ErrorResponse
// @Router       /resourceServers [get]
// @Security	BearerAuth
func (r *resourceServer) GetAllResourceServers(ctx *gin.Context) {
	var filtersParam db_pgnflt.PgnFltQueryParams
	err := ctx.BindQuery(&filtersParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		r.logger.Info(ctx, "invalid query params", zap.Error(err), zap.Any("query-params", ctx.Request.URL.Query()))
		_ = ctx.Error(err)
		return
	}

	resourceServers, metaData, err := r.resourceServerModule.GetAllResourceServers(ctx, filtersParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, resourceServers, metaData)
}
