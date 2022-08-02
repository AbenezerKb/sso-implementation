package initiator

import (
	"context"
	"fmt"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sso/internal/handler/middleware"
	"sso/platform/logger"
	"syscall"
	"time"
)

func Initiate() {
	log := logger.New(InitLogger())
	log.Info(context.Background(), "logger initialized")

	log.Info(context.Background(), "initializing config")
	InitConfig("config", "config", log)
	log.Info(context.Background(), "config initialized")

	log.Info(context.Background(), "initializing database")
	db := InitDB(viper.GetString("database.url"), log)
	log.Info(context.Background(), "database initialized")

	log.Info(context.Background(), "initializing cache")
	cache := InitCache(viper.GetString("redis.url"), log)
	log.Info(context.Background(), "cache initialized")

	log.Info(context.Background(), "initializing persistence layer")
	persistence := InitPersistence(db, log)
	log.Info(context.Background(), "persistence layer initialized")

	log.Info(context.Background(), "initializing module")
	module := InitModule(persistence, cache, log)
	log.Info(context.Background(), "module initialized")

	log.Info(context.Background(), "initializing handler")
	handler := InitHandler(module, log)
	log.Info(context.Background(), "handler initialized")

	log.Info(context.Background(), "initializing server")
	server := gin.New()
	server.Use(middleware.GinLogger(log))
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin.recovery"), true))
	log.Info(context.Background(), "server initialized")

	log.Info(context.Background(), "initializing router")
	v1 := server.Group("/v1")
	InitRouter(server, v1, handler, module, log)
	log.Info(context.Background(), "router initialized")

	srv := &http.Server{
		Addr:    viper.GetString("server.host") + ":" + viper.GetString("server.port"),
		Handler: server,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		log.Info(context.Background(), "server started",
			zap.String("host", viper.GetString("server.host")),
			zap.Int("port", viper.GetInt("server.port")))
		log.Info(context.Background(), fmt.Sprintf("server stopped with error %v", srv.ListenAndServe()))
	}()
	sig := <-quit
	log.Info(context.Background(), fmt.Sprintf("server shutting down with signal %v", sig))
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("server.timeout")*time.Second)
	defer cancel()

	log.Info(ctx, "shutting down server")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("error while shutting down server: %v", err))
	} else {
		log.Info(context.Background(), "server shutdown complete")
	}
}
