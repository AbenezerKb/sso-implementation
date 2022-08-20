package utils

import (
	"context"
	"fmt"
	"math/rand"
	"sso/platform/logger"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

	rand.Seed(time.Now().Unix())
	randString := make([]byte, length)
	for i := range randString {
		randString[i] = str[rand.Int63()%int64(len(str))]
	}
	return string(randString)
}

func GenerateTimeStampedRandomString(length int, includeSpecial bool) string {
	return fmt.Sprintf("%s%d", GenerateRandomString(length, includeSpecial), time.Now().Unix())
}

func ArrayToString(array []string) string {
	return strings.Join(array, " ")
}

func StringToArray(str string) []string {
	return strings.Split(str, " ")
}
