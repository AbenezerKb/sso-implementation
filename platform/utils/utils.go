package utils

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"sso/platform/logger"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
const specialBytes = `!@#$%^&*()_+-=;':"[]{},.<>`

func HashAndSalt(ctx context.Context, pwd []byte, logger logger.Logger) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		logger.Error(ctx, "could not hash password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

func GenerateRandomString(length int, includeSpecial bool) string {
	str := letterBytes
	if includeSpecial {
		str += specialBytes
	}
	randString := make([]byte, length)
	for i := range randString {
		randString[i] = str[rand.Int63()%int64(len(str))]
	}
	return string(randString)
}

func ArrayToString(array []string) string {
	return strings.Join(array, " ")
}

func StringToArray(str string) []string {
	return strings.Split(str, " ")
}
