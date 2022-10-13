package self

import (
	"context"
	"fmt"
	"net/http"
	"sso/internal/constant/model/dto"
	"sso/platform"
	"sso/platform/utils"
)

type identityProvider struct {
}

func Init() platform.IdentityProvider {
	return &identityProvider{}
}
func (i *identityProvider) GetAccessToken(ctx context.Context, endPoint, redirectURI, clientID, clientSecret, code string) (string, string, error) {
	if ctx.Err() != nil {
		return "", "", ctx.Err()
	}
	var bodyMap = map[string]interface{}{
		"code":         code,
		"redirect_uri": redirectURI,
		"grant_type":   "authorization_code",
	}
	var response struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	res, err := utils.DoRequest(http.MethodPost, endPoint, bodyMap, &response, func(req *http.Request) error {
		req.SetBasicAuth(clientID, clientSecret)
		req.Header.Add("Content-Type", "application/json")
		return nil
	})
	if err != nil {
		return "", "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("identity provider responded with status: %d", res.StatusCode)
	}

	return response.Data.AccessToken, response.Data.RefreshToken, nil
}

func (i *identityProvider) GetUserInfo(ctx context.Context, endPoint, accessToken string) (dto.UserInfo, error) {
	if ctx.Err() != nil {
		return dto.UserInfo{}, ctx.Err()
	}
	var response struct {
		Data dto.UserInfo
	}
	res, err := utils.DoRequest(http.MethodGet, endPoint, nil, &response, func(req *http.Request) error {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessToken)
		return nil
	})
	if err != nil {
		return dto.UserInfo{}, err
	}
	if res.StatusCode != http.StatusOK {
		return dto.UserInfo{}, fmt.Errorf("identity provider responded with status: %d", res.StatusCode)
	}

	return response.Data, nil
}
