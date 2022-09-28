package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sso/initiator"
	"sso/internal/constant"
	"sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/middleware"
	"sso/platform/logger"
	"sso/platform/utils"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/segmentio/kafka-go"

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
	response struct {
		OK   bool              `json:"ok"`
		Data dto.TokenResponse `json:"data"`
	}
	AccessToken   string
	RefreshToken  string
	enforcer      *casbin.Enforcer
	Logger        logger.Logger
	Conn          *pgxpool.Pool
	PlatformLayer initiator.PlatformLayer
	CacheLayer    initiator.CacheLayer
	KafkaConn     *kafka.Conn
	KafkaReader   *kafka.Reader
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
	persistence := initiator.InitPersistence(sqlConn, db.NewPhoneSwap(pgxConn), log)
	log.Info(context.Background(), "persistence layer initialized")

	log.Info(context.Background(), "initializing cache layer")
	cacheLayer := initiator.InitMockCacheLayer(cache, viper.GetDuration("redis.otp_expire_time"), "123455", log, initiator.CacheOptions{
		OTPExpireTime:      viper.GetDuration("redis.otp_expire_time"),
		SessionExpireTime:  viper.GetDuration("redis.session_expire_time"),
		ConsentExpireTime:  viper.GetDuration("redis.consent_expire_time"),
		AuthCodeExpireTime: viper.GetDuration("redis.authcode_expire_time"),
	})
	log.Info(context.Background(), "cache layer initialized")

	log.Info(context.Background(), "initializing platform layer")
	platformLayer := initiator.InitMockPlatformLayer(log, path+viper.GetString("private_key"), path+viper.GetString("public_key"))
	log.Info(context.Background(), "platform layer initialized")

	log.Info(context.Background(), "initializing state")
	state := initiator.InitState(log)
	log.Info(context.Background(), "state initialized")

	log.Info(context.Background(), "initializing module")
	module := initiator.InitMockModule(persistence, cacheLayer, path+viper.GetString("private_key"), platformLayer, log, enforcer, state)
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
	initiator.InitRouter(server, v1, handler, module, log, enforcer, platformLayer)
	log.Info(context.Background(), "router initialized")
	kafkaConn := kafkaConn(viper.GetString("kafka.url"), viper.GetString("kafka.topic"))

	kafkaReader := kafkaReader(viper.GetString("kafka.url"), viper.GetString("kafka.topic"), viper.GetString("kafka.group_id"))
	return TestInstance{
		Server:        server,
		DB:            sqlConn,
		Redis:         cache,
		Module:        module,
		enforcer:      enforcer,
		Logger:        log,
		Conn:          pgxConn,
		PlatformLayer: platformLayer,
		CacheLayer:    cacheLayer,
		KafkaConn:     kafkaConn,
		KafkaReader:   kafkaReader,
	}
}
func (t *TestInstance) Authenticate(credentials *godog.Table) (db.User, error) {
	// read email and password from table
	apiTest := src.ApiTest{
		URL:    "/v1/login",
		Method: http.MethodPost,
	}
	email, err := apiTest.ReadCellString(credentials, "email")
	if err != nil {
		return db.User{}, err
	}
	password, err := apiTest.ReadCellString(credentials, "password")
	if err != nil {
		return db.User{}, err
	}
	hash, err := utils.HashAndSalt(context.Background(), []byte(password), t.Logger)
	if err != nil {
		return db.User{}, err
	}
	user, err := t.DB.CreateUser(context.Background(), db.CreateUserParams{
		Email:    utils.StringOrNull(email),
		Password: hash,
	})
	if err != nil {
		return db.User{}, err
	}

	apiTest.InitializeServer(t.Server)
	apiTest.SetHeader("Content-Type", "application/json")
	body, err := apiTest.ReadRow(credentials, nil, false)
	if err != nil {
		return db.User{}, err
	}

	apiTest.Body = body
	apiTest.SendRequest()

	// if err := apiTest.AssertStatusCode(http.StatusBadRequest); err != nil {
	// 	return err
	// }

	err = json.Unmarshal(apiTest.ResponseBody, &t.response)
	if err != nil {
		return db.User{}, err
	}

	t.AccessToken = t.response.Data.AccessToken
	t.RefreshToken = t.response.Data.RefreshToken
	return user, nil
}

func (t *TestInstance) GrantRoleForUser(userID string, role *godog.Table) error {
	test := src.ApiTest{}
	permission, err := test.ReadCellString(role, "role")
	if err != nil {
		return err
	}
	exists, err := t.enforcer.HasRoleForUser(userID, permission, constant.User)
	if err != nil {
		return err
	}
	if !exists {
		_, err := t.enforcer.AddRoleForUser(userID, permission, constant.User)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TestInstance) AuthenticateWithParam(credentials dto.User) (db.User, error) {
	apiTest := src.ApiTest{
		URL:    "/v1/login",
		Method: http.MethodPost,
	}

	apiTest.InitializeServer(t.Server)
	apiTest.SetHeader("Content-Type", "application/json")
	apiTest.SetBodyMap(map[string]interface{}{
		"password": credentials.Password,
		"email":    credentials.Email,
	})

	var err error
	credentials.Password, err = utils.HashAndSalt(context.Background(), []byte(credentials.Password), t.Logger)
	if err != nil {
		return db.User{}, err
	}
	user, err := t.DB.CreateUser(context.Background(), db.CreateUserParams{
		FirstName:      credentials.FirstName,
		MiddleName:     credentials.MiddleName,
		LastName:       credentials.LastName,
		Email:          sql.NullString{String: credentials.Email, Valid: true},
		Phone:          credentials.Phone,
		Password:       credentials.Password,
		UserName:       credentials.UserName,
		Gender:         credentials.Gender,
		ProfilePicture: sql.NullString{String: credentials.ProfilePicture, Valid: true},
	})
	if err != nil {
		return db.User{}, err
	}

	apiTest.SendRequest()
	err = json.Unmarshal(apiTest.ResponseBody, &t.response)
	if err != nil {
		return db.User{}, err
	}

	t.AccessToken = t.response.Data.AccessToken
	t.RefreshToken = t.response.Data.RefreshToken
	return user, nil
}

func kafkaConn(address, topic string) *kafka.Conn {
	KafkaConn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, 0)
	if err != nil {
	}

	return KafkaConn
}

func kafkaReader(address, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(address, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
}
