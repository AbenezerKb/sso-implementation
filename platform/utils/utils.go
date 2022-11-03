package utils

import (
	"context"
	"crypto"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sso/platform/logger"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
const specialBytes = `!@#$%^&*:.`

type CookieOptions struct {
	Path, Domain     string
	MaxAge           int
	Secure, HttpOnly bool
	SameSite         int
}

func HashAndSalt(ctx context.Context, pwd []byte, logger logger.Logger) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		logger.Error(ctx, "could not hash password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}
func CompareHashAndPassword(hashedPwd, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPassword))
	return err == nil
}

// GenerateRandomString generates a random string with the specified length
//
// as of oct-25-2022 includeSpecial has no effect on the output
func GenerateRandomString(length int, includeSpecial bool) string {
	str := letterBytes

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
	return fmt.Sprintf("%x.%s", hash.Sum(nil), salt)
}

func GenerateRedirectString(uri *url.URL, queries map[string]string) string {
	query := uri.Query()
	for k, v := range queries {
		query.Set(k, v)
	}
	uri.RawQuery = query.Encode()
	return uri.String()
}

func SaveMultiPartFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)

	return err
}

func SetRefreshTokenCookie(ctx *gin.Context, value string, options CookieOptions) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "ab_fen",
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: http.SameSite(options.SameSite),
	})
}

func RemoveRefreshTokenCookie(ctx *gin.Context, options CookieOptions) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "ab_fen",
		Value:    "",
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   -1,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: http.SameSite(options.SameSite),
	})
}

func SetOPBSCookie(ctx *gin.Context, value string, options CookieOptions) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "opbs",
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: http.SameSite(options.SameSite),
	})
}
