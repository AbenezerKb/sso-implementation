package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sso/internal/constant/model/db"
	"sso/internal/handler/middleware"
	"sso/platform/logger"
	"syscall"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Initiate
// @title           RidePLUS SSO API
// @version         0.1
// @description     This is the RidePLUS sso api.
//
// @contact.name   2F Capital Support Email
// @contact.url    http://www.2fcapital.com
// @contact.email  info@1f-capital.com
//
// @host 206.189.54.235:8000
// @BasePath  /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.basic BasicAuth
func Initiate() {
	log := logger.New(InitLogger())
	log.Info(context.Background(), "logger initialized")

	log.Info(context.Background(), "initializing config")
	configName := "config"
	if name := os.Getenv("CONFIG_NAME"); name != "" {
		configName = name
		log.Info(context.Background(), fmt.Sprintf("config name is set to %s", configName))
	} else {
		log.Info(context.Background(), "using default config name 'config'")
	}
	InitConfig(configName, "config", log)
	log.Info(context.Background(), "config initialized")

	log.Info(context.Background(), "initializing database")
	pgxConn := InitDB(viper.GetString("database.url"), log)
	log.Info(context.Background(), "database initialized")

	log.Info(context.Background(), "initializing migration")
	m := InitiateMigration(viper.GetString("migration.path"), viper.GetString("database.url"), log)
	UpMigration(m, log)
	log.Info(context.Background(), "migration initialized")

	log.Info(context.Background(), "initializing cache")
	cache := InitCache(viper.GetString("redis.url"), log)
	log.Info(context.Background(), "cache initialized")

	log.Info(context.Background(), "initializing casbin enforcer")
	enforcer := InitEnforcer(viper.GetString("casbin.path"), pgxConn, log)
	log.Info(context.Background(), "casbin enforcer initialized")

	log.Info(context.Background(), "initializing persistence layer")
	persistence := InitPersistence(db.New(pgxConn), log)
	log.Info(context.Background(), "persistence layer initialized")

	log.Info(context.Background(), "initializing cache layer")
	cacheLayer := InitCacheLayer(cache, CacheOptions{
		OTPExpireTime:      viper.GetDuration("redis.otp_expire_time"),
		SessionExpireTime:  viper.GetDuration("redis.session_expire_time"),
		ConsentExpireTime:  viper.GetDuration("redis.consent_expire_time"),
		AuthCodeExpireTime: viper.GetDuration("redis.authcode_expire_time"),
	}, log)
	log.Info(context.Background(), "cache layer initialized")

	log.Info(context.Background(), "initializing platform layer")
	platformLayer := InitPlatformLayer(log, viper.GetString("private_key"), viper.GetString("public_key"))
	log.Info(context.Background(), "platform layer initialized")

	log.Info(context.Background(), "initializing state")
	state := InitState(log)
	log.Info(context.Background(), "state initialized")

	log.Info(context.Background(), "initializing module")
	module := InitModule(persistence, cacheLayer, viper.GetString("private_key"), platformLayer, log, enforcer, state)
	log.Info(context.Background(), "module initialized")

	log.Info(context.Background(), "initializing handler")
	handler := InitHandler(module, log)
	log.Info(context.Background(), "handler initialized")

	log.Info(context.Background(), "initializing server")
	server := gin.New()
	server.Use(middleware.GinLogger(log.Named("gin")))
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin.recovery"), true))
	server.Use(middleware.ErrorHandler())
	if viper.GetBool("dev") {
		server.Use(InitCORS())
	}
	log.Info(context.Background(), "server initialized")

	log.Info(context.Background(), "initializing metrics route")
	InitMetricsRoute(server, log)
	log.Info(context.Background(), "metrics route initialized")

	log.Info(context.Background(), "initializing router")
	v1 := server.Group("/v1")
	InitRouter(server, v1, handler, module, log, enforcer, platformLayer)
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
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("server.timeout"))
	defer cancel()

	log.Info(ctx, "shutting down server")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("error while shutting down server: %v", err))
	} else {
		log.Info(context.Background(), "server shutdown complete")
	}
}
