package test

import (
	"context"
	"fmt"
	"os"
	"sso/initiator"
	"sso/internal/constant/model/db"
	"sso/internal/handler/middleware"
	"sso/platform/logger"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type TestInstance struct {
	server *gin.Engine
	db     *db.Queries
}

var Instance TestInstance = TestInstance{}

func Initiate(path string) {
	log := logger.New(initiator.InitLogger())
	log.Info(context.Background(), "logger initialized")

	configName := "config"
	if name := os.Getenv("CONFIG_NAME"); name != "" {
		configName = name
		log.Info(context.Background(), fmt.Sprintf("config name is set to %s", configName))
	} else {
		log.Info(context.Background(), "using default config name 'config'")
	}
	log.Info(context.Background(), "initializing config")
	initiator.InitConfig(configName, path, log)
	log.Info(context.Background(), "config initialized")

	log.Info(context.Background(), "initializing database")
	db := initiator.InitDB(viper.GetString("database.url"), log)
	log.Info(context.Background(), "database initialized")

	log.Info(context.Background(), "initializing cache")
	cache := initiator.InitCache(viper.GetString("redis.url"), log)
	log.Info(context.Background(), "cache initialized")

	log.Info(context.Background(), "initializing persistence layer")
	persistence := initiator.InitPersistence(db, log)
	log.Info(context.Background(), "persistence layer initialized")

	log.Info(context.Background(), "initializing module")
	module := initiator.InitModule(persistence, cache, log)
	log.Info(context.Background(), "module initialized")

	log.Info(context.Background(), "initializing handler")
	handler := initiator.InitHandler(module, log)
	log.Info(context.Background(), "handler initialized")

	log.Info(context.Background(), "initializing server")
	server := gin.New()
	server.Use(middleware.GinLogger(log))
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin.recovery"), true))
	log.Info(context.Background(), "server initialized")

	log.Info(context.Background(), "initializing router")
	v1 := server.Group("/v1")
	initiator.InitRouter(server, v1, handler, module, log)
	log.Info(context.Background(), "router initialized")

	Instance = TestInstance{server, db}
	// return server, db
}

func GetServer(path string) (*gin.Engine, *db.Queries) {
	if Instance.server == nil {
		Initiate(path)
	}

	return Instance.server, Instance.db
}

type SqlcDB struct {
	DB *db.Queries
}

func (sq *SqlcDB) Feed(interface{}) error {
	return nil
}

func (sq *SqlcDB) Starve(interface{}) error {
	return nil
}
