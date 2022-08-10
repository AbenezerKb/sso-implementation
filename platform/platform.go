package platform

import "context"

type SMSConfig struct {
	UserName  string
	Password  string
	Sender    string
	DLRMask   string
	DCS       string
	DLRURL    string
	Server    string
	Templates map[string]string
	Type      string
	APIKey    string
}

type SMSClient interface {
	SendSMS(ctx context.Context, to, text string) error
	SendSMSWithTemplate(ctx context.Context, to, templateName string, values ...interface{}) error
}
