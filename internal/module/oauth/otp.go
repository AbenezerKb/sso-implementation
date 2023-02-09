package oauth

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

	"sso/internal/constant/errors"

	"github.com/dongri/phonenumber"
	"go.uber.org/zap"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func (o *oauth) GenerateOTP(ctx context.Context, phone string, max int) (string, error) {
	// if there is an existing otp, use that instead of generating a new one
	otp, err := o.otpCache.GetOTP(ctx, phone)
	if err != nil {
		return otp, err
	}
	if otp != "" {
		return otp, nil
	}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil {
		err := errors.ErrOTPGenerate.Wrap(err, "failed to generate otp")
		o.logger.Error(ctx, "failed to generate otp", zap.Error(err))
		return "", err
	}

	if n != max {
		err := errors.ErrOTPGenerate.New(fmt.Sprintf("incorrect otp with length %d was generated", n))
		o.logger.Error(ctx, "failed to generate otp", zap.Error(err))
		return "", err
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

func (o *oauth) RequestOTP(ctx context.Context, phone string, rqType string) error {
	phone = phonenumber.Parse(phone, "ET")
	if phone == "" {
		err := errors.ErrInvalidUserInput.New("invalid phone number")
		o.logger.Info(ctx, "invalid phone number", zap.Error(err))
		return err
	}

	exists, err := o.oauthPersistence.UserByPhoneExists(ctx, phone)
	if err != nil {
		return err
	}
	if rqType == "signup" {
		if exists {
			err = errors.ErrDataExists.New("user with this phone already exists")
			o.logger.Info(ctx, "user with this phone already exists", zap.Error(err))
			return err
		}
	} else if rqType == "login" {
		if !exists {
			err = errors.ErrNoRecordFound.New("user with this phone does not exists")
			o.logger.Info(ctx, "user with this phone does not exists", zap.Error(err))
			return err
		}
		err := o.VerifyUserStatus(ctx, phone)
		if err != nil {
			return err
		}

	} else if rqType == "change" {
		if exists {
			err = errors.ErrDataExists.New("user with this phone already exists")
			o.logger.Info(ctx, "user with this phone already exists", zap.Error(err))
			return err
		}
	} else {
		err = errors.ErrInvalidUserInput.New("invalid request type")
		o.logger.Info(ctx, "invalid request type", zap.Error(err))
		return err
	}

	otp, err := o.GenerateOTP(ctx, phone, 6)
	if err != nil {
		return err
	}

	var defaultFound bool

	for _, v := range o.options.ExcludedPhones.Phones {
		if phone == v {
			otp = o.options.ExcludedPhones.DefaultOTP
			defaultFound = true

			break
		}
	}

	err = o.otpCache.SetOTP(ctx, phone, otp)
	if err != nil {
		return err
	}

	if !defaultFound || o.options.ExcludedPhones.SendSMS {
		err = o.smsClient.SendSMSWithTemplate(ctx, phone, "otp", otp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *oauth) VerifyOTP(ctx context.Context, phone string, otp string) error {
	otpFromCache, err := o.otpCache.GetOTP(ctx, phone)
	if err != nil {
		return err
	}
	if otpFromCache == "" {
		err := errors.ErrInvalidUserInput.New("invalid otp")
		o.logger.Info(ctx, "invalid otp", zap.Error(err))
		return err
	}
	if otpFromCache != otp {
		err = errors.ErrInvalidUserInput.New("invalid otp")
		o.logger.Info(ctx, "invalid otp", zap.Error(err))
		return err
	}

	return o.otpCache.DeleteOTP(ctx, phone)
}
