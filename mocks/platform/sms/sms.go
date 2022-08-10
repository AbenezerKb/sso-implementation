package sms

import (
	"context"
	"sso/platform"
	"sso/platform/logger"
)

type mockSMSClient struct {
	smsConfig platform.SMSConfig
	logger    logger.Logger
}

func InitMockSMS(smsConfig platform.SMSConfig, logger logger.Logger) platform.SMSClient {
	return &mockSMSClient{
		smsConfig: smsConfig,
		logger:    logger,
	}
}

func (m *mockSMSClient) SendSMS(ctx context.Context, to, text string) error {
	return nil
}

func (m *mockSMSClient) SendSMSWithTemplate(ctx context.Context, to, templateName string, values ...interface{}) error {
	return nil
}
