package oauth

import (
	"context"
	"io/ioutil"
	"log"
	"sso/internal/constant/model/dto"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (o *oauth) GenerateAccessToken(ctx context.Context, user *dto.User) (string, error) {
	keyFile, err := ioutil.ReadFile("./privatekey.pem")
	if err != nil {
		log.Fatal("Error reading own private key : \n", err)
	}
	claims := dto.AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    "test",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return "", nil
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (o *oauth) GenerateRefreshToken(ctx context.Context, user *dto.User) (string, error) {
	keyFile, err := ioutil.ReadFile("./privatekey.pem")
	if err != nil {
		log.Fatal("Error reading own private key : \n", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return "", nil
	}

	cliams := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		Issuer:    "test",
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	rfToken, err := jwt.NewWithClaims(jwt.SigningMethodPS512, cliams).SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return rfToken, nil
}

func (o *oauth) GenerateIdToken(ctx context.Context, user *dto.User) (string, error) {
	keyFile, err := ioutil.ReadFile("./privatekey.pem")

	if err != nil {
		log.Fatal("Error reading own private key : \n", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return "", nil
	}

	claims := dto.IDTokenPayload{
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		Picture:     user.ProfilePicture,
		Email:       user.Email,
		PhoneNumber: user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(user.CreatedAt),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}
