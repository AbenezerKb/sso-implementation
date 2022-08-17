package initiator

import (
	"github.com/spf13/viper"
	sms2 "sso/mocks/platform/sms"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/sms"
)

type PlatformLayer struct {
	sms platform.SMSClient
}

func InitPlatformLayer(logger logger.Logger) PlatformLayer {
	return PlatformLayer{
		sms: sms.InitSMS(
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
	}
}

func InitMockPlatformLayer(logger logger.Logger) PlatformLayer {
	return PlatformLayer{
		sms: sms2.InitMockSMS(
			platform.SMSConfig{},
			logger.Named("sms-platform")),
	}
}
