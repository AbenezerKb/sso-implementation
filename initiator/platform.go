package initiator

import (
	"context"
	"crypto/rsa"
	"log"
	"os"

	"sso/internal/constant/model/dto"
	"sso/mocks/platform/identityProvider"
	sms2 "sso/mocks/platform/sms"
	"sso/platform"
	"sso/platform/asset"
	"sso/platform/identityProviders/self"
	kafka_consumer "sso/platform/kafka"
	"sso/platform/logger"
	"sso/platform/sms"
	"sso/platform/token"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type PlatformLayer struct {
	Sms    platform.SMSClient
	Token  platform.Token
	Kafka  kafka_consumer.Kafka
	SelfIP platform.IdentityProvider
	Asset  platform.Asset
}

func InitPlatformLayer(logger logger.Logger, privateKeyPath, publicKeyPath string, _ Persistence) PlatformLayer {

	return PlatformLayer{
		Sms: sms.InitSMS(
			platform.SMSConfig{
				UserName:  viper.GetString("sms.username"),
				Password:  viper.GetString("sms.password"),
				Server:    viper.GetString("sms.server"),
				Type:      viper.GetString("sms.type"),
				DCS:       viper.GetString("sms.dcs"),
				DLRMask:   viper.GetString("sms.dlrmask"),
				DLRURL:    viper.GetString("sms.dlrurl"),
				Sender:    viper.GetString("sms.sender"),
				Templates: viper.GetStringMapString("sms.templates"),
				APIKey:    viper.GetString("sms.api_key"),
			},
			logger.Named("sms-platform")),
		Token: token.JwtInit(logger.Named("token-platform"),
			privateKey(privateKeyPath),
			publicKey(publicKeyPath),
		),
		Kafka:  kafka_consumer.NewKafkaConnection(viper.GetString("kafka.url"), viper.GetString("kafka.group_id"), []string{viper.GetString("kafka.drivers_topic")}, viper.GetInt("kafka.max_read_bytes"), logger),
		SelfIP: self.Init(),
		Asset: asset.InitDigitalOceanAsset(logger.Named("asset-platform"),
			viper.GetString("digital_ocean.space.key"),
			viper.GetString("digital_ocean.space.secret"),
			viper.GetString("digital_ocean.space.url"),
			viper.GetString("digital_ocean.space.bucket"),
		),
	}
}

func InitMockPlatformLayer(logger logger.Logger, privateKeyPath, publicKeyPath string, _ Persistence) PlatformLayer {
	return PlatformLayer{
		Sms: sms2.InitMockSMS(
			platform.SMSConfig{},
			logger.Named("sms-platform")),
		Token: token.JwtInit(logger.Named("token-platform"),
			privateKey(privateKeyPath),
			publicKey(publicKeyPath),
		),
		// Kafka: kafka_consumer.NewKafkaConnection(viper.GetString("kafka.url"), viper.GetString("kafka.topic"), viper.GetString("kafka.group_id"), viper.GetInt("kafka.max_read_bytes"), logger),
		SelfIP: identityProvider.InitIP("some_id", "some_secret", "veryLegitCode", "legitAccessToken", dto.UserInfo{
			FirstName: "john",
			Email:     "john@gmail.com",
			Phone:     "0912131415",
		}),
		Asset: asset.Init(logger.Named("asset-platform"), "../../../../assets"),
	}
}

func privateKey(privateKeyPath string) *rsa.PrivateKey {
	keyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(context.Background(), "failed to read private key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		log.Fatal(context.Background(), "failed to parse private key", zap.Error(err))
	}
	return privateKey
}
func publicKey(publicKeyPath string) *rsa.PublicKey {
	certificate, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(context.Background(), "Error reading own certificate : \n", zap.Error(err))
	}
	ssoPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(certificate)
	if err != nil {
		log.Fatal(context.Background(), "Error parsing own certificate : \n", zap.Error(err))
	}
	return ssoPublicKey
}
