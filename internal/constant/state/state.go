package state

import (
	"context"
	"net/url"
	"sso/platform/logger"

	"github.com/spf13/viper"
)

const (
	ConsentKey  = "consent:%v"
	AuthCodeKey = "authcode:%v"
	//ConsentURL  = "https://www.google.com/"
	//ErrorURL    = "https://www.google.com/"
)

type URLs struct {
	ErrorURL   *url.URL
	ConsentURL *url.URL
	LogoutURL  *url.URL
}

func InitiateURLs(logger logger.Logger) URLs {
	errorURLString := viper.GetString("frontend.error_url")
	if errorURLString == "" {
		logger.Fatal(context.Background(), "unable to read frontend.error_url in viper")
	}
	ErrorURL, err := url.Parse(errorURLString)
	if err != nil {
		logger.Fatal(context.Background(), "unable to parse frontend.error_url")
	}
	consentURLString := viper.GetString("frontend.consent_url")
	if consentURLString == "" {
		logger.Fatal(context.Background(), "unable to read frontend.consent_url in viper")
	}
	ConsentURL, err := url.Parse(consentURLString)
	if err != nil {
		logger.Fatal(context.Background(), "unable to parse frontend.consent_url")
	}

	logoutURLString := viper.GetString("frontend.logout_url")
	if consentURLString == "" {
		logger.Fatal(context.Background(), "unable to read frontend.logout_url in viper")
	}
	LogoutURL, err := url.Parse(logoutURLString)
	if err != nil {
		logger.Fatal(context.Background(), "unable to parse frontend.logout_url")
	}
	return URLs{
		ErrorURL:   ErrorURL,
		ConsentURL: ConsentURL,
		LogoutURL:  LogoutURL,
	}
}
