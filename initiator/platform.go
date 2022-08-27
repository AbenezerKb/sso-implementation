package initiator

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"log"
	sms2 "sso/mocks/platform/sms"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/sms"
	"sso/platform/token"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type PlatformLayer struct {
	Sms   platform.SMSClient
	Token platform.Token
}

func InitPlatformLayer(logger logger.Logger, privateKeyPath, publicKeyPath string) PlatformLayer {
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
	}
}

func InitMockPlatformLayer(logger logger.Logger, privateKeyPath, publicKeyPath string) PlatformLayer {
	return PlatformLayer{
		Sms: sms2.InitMockSMS(
			platform.SMSConfig{},
			logger.Named("sms-platform")),
		Token: token.JwtInit(logger.Named("token-platform"),
			privateKey(privateKeyPath),
			publicKey(publicKeyPath),
		),
	}
}

func privateKey(privateKeyPath string) *rsa.PrivateKey {
	keyFile, err := ioutil.ReadFile(privateKeyPath)
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
	certificate, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(context.Background(), "Error reading own certificate : \n", zap.Error(err))
	}
	ssoPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(certificate)
	if err != nil {
		log.Fatal(context.Background(), "Error parsing own certificate : \n", zap.Error(err))
	}
	return ssoPublicKey
}
