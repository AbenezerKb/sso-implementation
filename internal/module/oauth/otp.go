package oauth

import (
	"context"
	"crypto/rand"
	"io"
	"sso/internal/constant/errors"

	"github.com/dongri/phonenumber"
	"go.uber.org/zap"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func (o *oauth) GenerateOTP(ctx context.Context, max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max || err != nil {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

func (o *oauth) RequestOTP(ctx context.Context, phone string, rqType string) error {
	phone = phonenumber.Parse(phone, "ET")

	exists, err := o.oauthPersistence.UserByPhoneExists(ctx, phone)
	if err != nil {
		return err
	}
	if rqType == "signup" && exists {
		err = errors.ErrDataExists.New("user with this phone already exists")
		o.logger.Error(ctx, "user with this phone already exists", zap.Error(err))
		return err

	} else if rqType == "login" {
		if !exists {
			err = errors.ErrNoRecordFound.New("user with this phone does not exists")
			o.logger.Error(ctx, "user with this phone does not exists", zap.Error(err))
			return err
		}
		err := o.VerifyUserStatus(ctx, phone)
		if err != nil {
			return err
		}

	} else if rqType != "signup" && rqType != "login" {
		err = errors.ErrInvalidUserInput.New("invalid request type")
		o.logger.Error(ctx, "invalid request type", zap.Error(err))
		return err
	}

	otp, err := o.GenerateOTP(ctx, 6)
	if err != nil {
		err = errors.ErrInternalServerError.New("error generating otp")
		o.logger.Error(ctx, "error generating otp", zap.Error(err))
		return err
	}
	err = o.otpCache.SetOTP(ctx, phone, otp)
	if err != nil {
		return err
	}
	err = o.SendSMS(ctx, phone, otp)
	if err != nil {
		return err
	}
	return nil
}

func (o *oauth) SendSMS(ctx context.Context, phone string, otp string) error {
	return nil
}

func (o *oauth) VerifyOTP(ctx context.Context, phone string, otp string) error {
	otpFromCache, err := o.otpCache.GetOTP(ctx, phone)
	if err != nil {
		return err
	}
	if otpFromCache != otp {
		err = errors.ErrInvalidUserInput.New("invalid credentials")
		o.logger.Info(ctx, "invalid otp", zap.Error(err))
		return err
	}
	return nil
}
