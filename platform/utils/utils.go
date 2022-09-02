package utils

import (
	"context"
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"sso/platform/logger"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
const specialBytes = `!@#$%^&*:.`

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
	_, _ = io.ReadAtLeast(rand.Reader, randString, length)
	for i := 0; i < len(randString); i++ {
		randString[i] = str[int(randString[i])%len(str)]
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

func ContainsValue(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func GenerateNewOPBS() string {
	return GenerateRandomString(100, false)
}

func CalculateSessionState(clientID, origin, opbs, salt string) string {
	hash := crypto.SHA256.New()
	hash.Write([]byte(fmt.Sprintf("%s %s %s %s", clientID, origin, opbs, salt)))
	return fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString(hash.Sum(nil)), salt)
}

func GenerateRedirectString(uri *url.URL, queries map[string]string) string {
	query := uri.Query()
	for k, v := range queries {
		query.Set(k, v)
	}
	uri.RawQuery = query.Encode()
	return uri.String()
}
