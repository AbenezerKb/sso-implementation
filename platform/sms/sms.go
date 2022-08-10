package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sso/internal/constant/errors"
	"sso/platform"
	"sso/platform/logger"
)

type smsModel struct {
	To   string `json:"to"`
	Text string `json:"text"`
}
type smsClient struct {
	platform.SMSConfig
	logger logger.Logger
}

func InitSMS(smsConfig platform.SMSConfig, logger logger.Logger) platform.SMSClient {
	return &smsClient{
		SMSConfig: smsConfig,
		logger:    logger,
	}
}

func (s *smsClient) SendSMS(ctx context.Context, to, text string) error {
	sms := smsModel{
		To:   to,
		Text: text,
	}
	js, err := json.Marshal(sms)
	if err != nil {
		err := errors.ErrSMSSend.Wrap(err, "couldn't send sms")
		s.logger.Error(ctx, "error while marshalling sms object", zap.Error(err))
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.Server, bytes.NewBuffer(js))
	if err != nil {
		err := errors.ErrSMSSend.Wrap(err, "couldn't send sms")
		s.logger.Error(ctx, "error while creating sms request", zap.Error(err))
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.UserName, s.Password)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		err := errors.ErrSMSSend.Wrap(err, "couldn't send sms")
		s.logger.Error(ctx, "error while sending sms request", zap.Error(err))
		return err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			s.logger.Error(ctx, "error while closing sms response", zap.Error(err))
		}
	}()

	if res.StatusCode != 202 {
		resBody, _ := ioutil.ReadAll(res.Body)
		err := errors.ErrSMSSend.Wrap(err, "couldn't send sms")
		s.logger.Error(ctx, fmt.Sprintf("error while sending sms. client responded with %v", string(resBody)), zap.Error(err))
		return err
	}

	return nil
}

func (s *smsClient) SendSMSWithTemplate(ctx context.Context, to, templateName string, values ...interface{}) error {
	template, ok := s.Templates[templateName]
	if !ok {
		s.logger.Warn(ctx, fmt.Sprintf("template %s not found. sending values only", templateName))
		template = "%v"
	}

	return s.SendSMS(ctx, to, fmt.Sprintf(template, values...))
}
