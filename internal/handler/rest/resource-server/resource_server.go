package resource_server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
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
// @Router       /resourceServer [post]
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
