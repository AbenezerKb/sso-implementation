package oauth2

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type oauth2 struct {
	logger            logger.Logger
	oauth2Persistence storage.OAuth2Persistence
	oauthPersistence  storage.OAuthPersistence
	consentCache      storage.ConsentCache
}

func InitOAuth2(logger logger.Logger, oauth2Persistence storage.OAuth2Persistence, oauthPersistence storage.OAuthPersistence, consentCache storage.ConsentCache) module.OAuth2Module {
	return &oauth2{
		logger:            logger,
		oauth2Persistence: oauth2Persistence,
		oauthPersistence:  oauthPersistence,
		consentCache:      consentCache,
	}
}

func (o *oauth2) Authorize(ctx context.Context, authRequestParma dto.AuthorizationRequestParam) (string, error) {
	if err := authRequestParma.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return "", err
	}
	client, err := o.oauth2Persistence.GetClient(ctx, authRequestParma.ClientID)
	if err != nil {
		return "", err
	}

	if !o.ContainsRedirectURL(client.RedirectURIs, authRequestParma.RedirectURI) {
		err := errors.ErrInvalidUserInput.New("invalid redirect uri")
		o.logger.Info(ctx, "invalid redirect uri", zap.Error(err))
		return "", err
	}

	scopes, err := o.oauth2Persistence.GetNamedScopes(ctx, strings.Split(authRequestParma.Scope, " ")...)
	if err != nil || len(scopes) == 0 {
		err := errors.ErrInvalidUserInput.New("invalid scope")
		o.logger.Info(ctx, "invalid scope", zap.Error(err))
		return "", err
	}

	consent := dto.Consent{
		ID:                        uuid.New().String(),
		AuthorizationRequestParam: authRequestParma,
	}
	if err := o.consentCache.SaveConsent(ctx, consent); err != nil {
		return "", err
	}

	return consent.ID, nil
}

// ContainsRedirectURL
func (o *oauth2) ContainsRedirectURL(redirectURIs []string, redirectURI string) bool {
	for _, ru := range redirectURIs {
		if ru == redirectURI {
			return true
		}
	}
	return false
}

func (o *oauth2) GetConsentByID(ctx context.Context, consentID string, id string) (dto.ConsentData, error) {
	consent, err := o.consentCache.GetConsent(ctx, consentID)
	if err != nil {
		err = errors.ErrNoRecordFound.Wrap(err, "consent not found")
		o.logger.Info(ctx, "consent not found", zap.Error(err))
		return dto.ConsentData{}, err
	}

	// get client
	client, err := o.oauth2Persistence.GetClient(ctx, consent.ClientID)
	if err != nil {
		return dto.ConsentData{}, err
	}
	// get scopes
	scopes, err := o.oauth2Persistence.GetNamedScopes(ctx, strings.Split(consent.Scope, " ")...)
	if err != nil {
		return dto.ConsentData{}, err
	}
	// get user

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "could not parse user id")
		o.logger.Error(ctx, "parse error", zap.Error(err))
		return dto.ConsentData{}, err
	}
	user, err := o.oauthPersistence.GetUserByID(ctx, userID)
	if err != nil {
		return dto.ConsentData{}, err
	}

	return dto.ConsentData{
		Consent: consent,
		Client:  client,
		Scopes:  scopes,
		User:    user,
	}, nil
}

func (o *oauth2) Approval(ctx context.Context, consentId string, accessRqResult string) (dto.Consent, error) {
	consent := dto.Consent{}

	if accessRqResult == "true" {
		var err error
		consent, err = o.consentCache.ChangeStatus(ctx, true, consentId)
		if err != nil {
			return dto.Consent{}, err
		}
	} else {
		var err error
		consent, err = o.consentCache.ChangeStatus(ctx, false, consentId)
		if err != nil {
			return dto.Consent{}, err
		}
	}
	return consent, nil
}

func (o *oauth2) IssueAuthCode(ctx context.Context, consent dto.Consent) (string, string, error) {
	authCode := dto.AuthCode{
		Code:        uuid.New().String(),
		Scope:       consent.Scope,
		RedirectURI: consent.AuthorizationRequestParam.RedirectURI,
		ClientID:    consent.ClientID,
		UserID:      consent.UserID,
	}
	if err := o.oauth2Persistence.SaveAuthCode(ctx, authCode); err != nil {
		return "", consent.State, err
	}
	return authCode.Code, consent.State, nil
}
