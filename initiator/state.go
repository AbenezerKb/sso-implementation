package initiator

import (
	"context"
	"net/url"

	"sso/internal/constant/state"
	"sso/platform/asset"
	"sso/platform/logger"

	"github.com/spf13/viper"
)

type State struct {
	URLs         state.URLs
	UploadParams state.UploadParams
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
	return State{
		URLs: state.URLs{
			ErrorURL:   errorURL,
			ConsentURL: consentURL,
			LogoutURL:  logoutURL,
		},
		UploadParams: asset.SetParams(logger, state.UploadParams{
			FileTypes: fileTypes,
		}),
	}
}
