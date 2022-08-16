package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sso/initiator"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/middleware"
	"sso/platform/logger"

	"github.com/casbin/casbin/v2"
	"github.com/cucumber/godog"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gitlab.com/2ftimeplc/2fbackend/bdd-testing-framework/src"
)

type TestInstance struct {
	Server   *gin.Engine
	DB       *db.Queries
	Redis    *redis.Client
	Module   initiator.Module
	User     dto.User
	response struct {
		OK   bool              `json:"ok"`
		Data dto.TokenResponse `json:"data"`
	}
	AccessToken string
	enforcer    *casbin.Enforcer
}

func Initiate(path string) TestInstance {
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
	initiator.InitConfig(configName, path+"config", log)
	log.Info(context.Background(), "config initialized")

	log.Info(context.Background(), "initializing database")
	pgxConn := initiator.InitDB(viper.GetString("database.url"), log)
	sqlConn := db.New(pgxConn)
	log.Info(context.Background(), "database initialized")

	log.Info(context.Background(), "initializing migration")
	m := initiator.InitiateMigration(path+viper.GetString("migration.path"), viper.GetString("database.url"), log)
	initiator.UpMigration(m, log)
	log.Info(context.Background(), "migration initialized")

	log.Info(context.Background(), "initializing casbin enforcer")
	enforcer := initiator.InitEnforcer(path+viper.GetString("casbin.path"), pgxConn, log)
	log.Info(context.Background(), "casbin enforcer initialized")

	log.Info(context.Background(), "initializing cache")
	cache := initiator.InitCache(viper.GetString("redis.url"), log)
	log.Info(context.Background(), "cache initialized")

	log.Info(context.Background(), "initializing persistence layer")
	persistence := initiator.InitPersistence(sqlConn, log)
	log.Info(context.Background(), "persistence layer initialized")

	log.Info(context.Background(), "initializing cache layer")
	cacheLayer := initiator.InitMockCacheLayer(cache, viper.GetDuration("redis.otp_expire_time"), "123455", log, initiator.CacheOptions{
		OTPExpireTime:     viper.GetDuration("redis.otp_expire_time"),
		SessionExpireTime: viper.GetDuration("redis.session_expire_time"),
		ConsentExpireTime: viper.GetDuration("redis.consent_expire_time"),
	})
	log.Info(context.Background(), "cache layer initialized")

	log.Info(context.Background(), "initializing platform layer")
	platformLayer := initiator.InitMockPlatformLayer(log)
	log.Info(context.Background(), "platform layer initialized")

	log.Info(context.Background(), "initializing module")
	module := initiator.InitModule(persistence, cacheLayer, path+viper.GetString("private_key"), platformLayer, log, enforcer)
	log.Info(context.Background(), "module initialized")

	log.Info(context.Background(), "initializing handler")
	handler := initiator.InitHandler(module, log)
	log.Info(context.Background(), "handler initialized")

	log.Info(context.Background(), "initializing server")
	server := gin.New()
	server.Use(middleware.GinLogger(log))
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin.recovery"), true))
	server.Use(middleware.ErrorHandler())
	log.Info(context.Background(), "server initialized")

	log.Info(context.Background(), "initializing metrics route")
	initiator.InitMetricsRoute(server, log)
	log.Info(context.Background(), "metrics route initialized")

	log.Info(context.Background(), "initializing router")
	v1 := server.Group("/v1")
	initiator.InitRouter(server, v1, handler, module, log, enforcer, path+viper.GetString("public_key"))
	log.Info(context.Background(), "router initialized")

	return TestInstance{
		Server:   server,
		DB:       sqlConn,
		Redis:    cache,
		Module:   module,
		enforcer: enforcer,
	}
}
func (t *TestInstance) Authenicate(creadentials *godog.Table) error {
	// read email and password from table
	apiTest := src.ApiTest{
		URL:    "/v1/login",
		Method: http.MethodPost,
	}
	email, err := apiTest.ReadCellString(creadentials, "email")
	if err != nil {
		return err
	}
	password, err := apiTest.ReadCellString(creadentials, "password")
	if err != nil {
		return err
	}
	hash, err := t.Module.OAuthModule.HashAndSalt(context.Background(), []byte(password))
	if err != nil {
		return err
	}
	userData, err := t.DB.CreateUser(context.Background(), db.CreateUserParams{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Password: hash,
	})
	if err != nil {
		return err
	}

	t.User.ID = userData.ID
	apiTest.InitializeServer(t.Server)
	apiTest.SetHeader("Content-Type", "application/json")
	body, err := apiTest.ReadRow(creadentials, nil, false)
	if err != nil {
		return err
	}

	apiTest.Body = body
	apiTest.SendRequest()

	// if err := apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
	// 	return err
	// }

	err = json.Unmarshal(apiTest.ResponseBody, &t.response)
	if err != nil {
		return err
	}

	permission, err := apiTest.ReadCellString(creadentials, "role")
	if err != nil {
		return err
	}
	t.enforcer.AddRoleForUser(t.User.ID.String(), permission)
	t.AccessToken = t.response.Data.AccessToken
	return nil
}
