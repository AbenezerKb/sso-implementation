package initiator

import (
	"context"
	"net/url"

	"sso/internal/constant/state"
	"sso/platform/asset"
	"sso/platform/logger"

	"github.com/dongri/phonenumber"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type State struct {
	URLs           state.URLs
	UploadParams   state.UploadParams
	ExcludedPhones state.ExcludedPhones
}

func InitState(logger logger.Logger) State {
	assets := GetMapSlice("assets")
	fileTypes := make([]state.FileType, 0, len(assets))

	for _, v := range assets {
		var fileType state.FileType

		fileType.SetValues(v)
		fileTypes = append(fileTypes, fileType)
	}

	errorURLString := viper.GetString("frontend.error_url")
	if errorURLString == "" {
		logger.Fatal(context.Background(), "unable to read frontend.error_url in viper")
	}
	errorURL, err := url.Parse(errorURLString)
	if err != nil {
		logger.Fatal(context.Background(), "unable to parse frontend.error_url")
	}
	consentURLString := viper.GetString("frontend.consent_url")
	if consentURLString == "" {
		logger.Fatal(context.Background(), "unable to read frontend.consent_url in viper")
	}
	consentURL, err := url.Parse(consentURLString)
	if err != nil {
		logger.Fatal(context.Background(), "unable to parse frontend.consent_url")
	}

	logoutURLString := viper.GetString("frontend.logout_url")
	if consentURLString == "" {
		logger.Fatal(context.Background(), "unable to read frontend.logout_url in viper")
	}
	logoutURL, err := url.Parse(logoutURLString)
	if err != nil {
		logger.Fatal(context.Background(), "unable to parse frontend.logout_url")
	}

	phones := viper.GetStringSlice("excluded_phones.phones")
	defaultOTP := viper.GetString("excluded_phones.default_otp")
	sendSMS := viper.GetBool("excluded_phones.send_sms")

	logger.Info(context.Background(), "using default otp",
		zap.String("default-otp", defaultOTP),
		zap.Bool("send-sms", sendSMS))

	for k, v := range phones {
		phone := phonenumber.Parse(v, "ET")
		if phone == "" {
			logger.Fatal(context.Background(),
				"invalid phone number for excluded phones", zap.String("phone", v))
		}

		phones[k] = phone
	}

	return State{
		URLs: state.URLs{
			ErrorURL:   errorURL,
			ConsentURL: consentURL,
			LogoutURL:  logoutURL,
		},
		UploadParams: asset.SetParams(logger, state.UploadParams{
			FileTypes: fileTypes,
		}),
		ExcludedPhones: state.ExcludedPhones{
			DefaultOTP: defaultOTP,
			Phones:     phones,
			SendSMS:    sendSMS,
		},
	}
}
