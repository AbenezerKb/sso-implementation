package identityProvider

import (
	"context"
	"fmt"
	"sso/internal/constant/model/dto"
	"sso/platform"
)

type identityProvider struct {
	clientID, clientSecret, legitCode, accessToken string
	user                                           dto.UserInfo
}

func InitIP(clientID, clientSecret, legitCode, accessToken string, user dto.UserInfo) platform.IdentityProvider {
	return &identityProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		accessToken:  accessToken,
		legitCode:    legitCode,
		user:         user,
	}
}
func SetUserForProvider(user dto.UserInfo, provider *platform.IdentityProvider) error {
	p, ok := (*provider).(*identityProvider)
	if !ok {
		return fmt.Errorf("invalid provider")
	}
	p.user = user
	return nil
}
func (i *identityProvider) GetAccessToken(ctx context.Context, endPoint, redirectURI, clientID, clientSecret, code string) (string, string, error) {
	if clientID == i.clientID && clientSecret == i.clientSecret {
		if code == i.legitCode {
			return i.accessToken, "", nil
		} else {
			return "", "", fmt.Errorf("invalid code")
		}
	}
	return "", "", fmt.Errorf("unauthorized")
}

func (i *identityProvider) GetUserInfo(ctx context.Context, endPoint, accessToken string) (dto.UserInfo, error) {
	if accessToken == accessToken {
		return i.user, nil
	}

	return dto.UserInfo{}, fmt.Errorf("invalid access token")
}
