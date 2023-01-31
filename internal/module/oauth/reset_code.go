package oauth

import (
	"context"
	"crypto/rand"
	"io"

	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"

	"github.com/dongri/phonenumber"
	"go.uber.org/zap"
)

func generateResetCode() string {
	str := "abcdefghijklmnopqrstuvwxyz0123456789"

	randString := make([]byte, 10)
	_, _ = io.ReadAtLeast(rand.Reader, randString, 10) //nolint:errcheck // since length = len(randString)

	for i := 0; i < len(randString); i++ {
		randString[i] = str[int(randString[i])%len(str)]
	}

	return string(randString)
}

func (o *oauth) generateResetCode(ctx context.Context, phone string) (string, error) {
	// if there is an existing resetCode, use that instead of generating a new one
	resetCode, err := o.resetCodeCache.GetResetCode(ctx, phone)
	if err == nil {
		return resetCode, err
	}
	if resetCode != "" {
		return resetCode, nil
	}

	return generateResetCode(), nil
}

func (o *oauth) RequestResetCode(ctx context.Context, phone string) error {
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
	if !exists {
		err = errors.ErrNoRecordFound.New("user with this phone does not exists")
		o.logger.Info(ctx, "user with this phone does not exists", zap.Error(err), zap.String("phone", phone))
		return nil // hide the error for security
	}
	err = o.VerifyUserStatus(ctx, phone)
	if err != nil {
		return err
	}

	resetCode, err := o.generateResetCode(ctx, phone)
	if err != nil {
		return err
	}
	err = o.resetCodeCache.SaveResetCode(ctx, phone, resetCode)
	if err != nil {
		return err
	}
	err = o.smsClient.SendSMSWithTemplate(ctx, phone, "reset_code", resetCode)
	if err != nil {
		return err
	}
	return nil
}

func (o *oauth) verifyResetCode(ctx context.Context, phone string, resetCode string) error {
	resetCodeFromCache, err := o.resetCodeCache.GetResetCode(ctx, phone)
	if err != nil {
		return err
	}
	if resetCodeFromCache == "" {
		err := errors.ErrInvalidUserInput.New("invalid reset code")
		o.logger.Info(ctx, "invalid reset code", zap.Error(err))
		return err
	}
	if resetCodeFromCache != resetCode {
		err = errors.ErrInvalidUserInput.New("invalid reset code")
		o.logger.Info(ctx, "invalid reset code", zap.Error(err))
		return err
	}

	return o.resetCodeCache.DeleteResetCode(ctx, phone)
}

func (o *oauth) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error {
	if err := request.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input")

		return err
	}

	request.Phone = phonenumber.Parse(request.Phone, "ET")

	// check code
	validCode, err := o.resetCodeCache.GetResetCode(ctx, request.Phone)
	if err != nil {
		return err
	}

	if validCode != request.ResetCode {
		err := errors.ErrInvalidUserInput.New("invalid reset code")
		o.logger.Info(ctx, "invalid reset code was tried", zap.String("reset-code", request.ResetCode))

		return err
	}

	// change password
	passwordHash, err := utils.HashAndSalt(ctx, []byte(request.Password), o.logger)
	if err != nil {
		return err
	}
	err = o.resetCodeCache.DeleteResetCode(ctx, request.Phone)
	if err != nil {
		return err
	}

	return o.oauthPersistence.ChangeUserPassword(ctx, request.Phone, passwordHash)
}
