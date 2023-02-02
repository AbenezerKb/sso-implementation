package oauth

import (
	"context"
	"crypto/rand"
	"io"

	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/joomcode/errorx"
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

func (o *oauth) generateResetCode(ctx context.Context, email string) (string, error) {
	// if there is an existing resetCode, use that instead of generating a new one
	resetCode, err := o.resetCodeCache.GetResetCode(ctx, email)
	if err == nil {
		return resetCode, err
	}
	if resetCode != "" {
		return resetCode, nil
	}

	return generateResetCode(), nil
}

func (o *oauth) RequestResetCode(ctx context.Context, email string) error {
	if err := validation.Validate(email,
		validation.Required.Error("email is required"),
		is.Email.Error("invalid email")); err != nil {
		err := errors.ErrInvalidUserInput.New("invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))

		return err
	}

	user, err := o.oauthPersistence.GetUserByEmail(ctx, email)
	if err != nil {
		if errorx.IsOfType(err, errors.ErrNoRecordFound) {
			return nil // hide the error for security
		}

		return err
	}

	err = o.VerifyUserStatus(ctx, user.Phone)
	if err != nil {
		return err
	}

	resetCode, err := o.generateResetCode(ctx, user.Email)
	if err != nil {
		return err
	}
	err = o.resetCodeCache.SaveResetCode(ctx, user.Email, resetCode)
	if err != nil {
		return err
	}
	err = o.smsClient.SendSMSWithTemplate(ctx, user.Phone, "reset_code", resetCode)
	if err != nil {
		return err
	}
	return nil
}

func (o *oauth) verifyResetCode(ctx context.Context, email string, resetCode string) error {
	resetCodeFromCache, err := o.resetCodeCache.GetResetCode(ctx, email)
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

	return o.resetCodeCache.DeleteResetCode(ctx, email)
}

func (o *oauth) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error {
	if err := request.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input")

		return err
	}

	// check code
	validCode, err := o.resetCodeCache.GetResetCode(ctx, request.Email)
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
	err = o.resetCodeCache.DeleteResetCode(ctx, request.Email)
	if err != nil {
		return err
	}

	return o.oauthPersistence.ChangeUserPassword(ctx, request.Email, passwordHash)
}
