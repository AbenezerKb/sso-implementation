package oauth2

import (
	"context"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Options struct {
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
}

func SetOptions(options Options) Options {
	if options.AccessTokenExpireTime == 0 {
		options.AccessTokenExpireTime = time.Minute * 10
	}
	if options.RefreshTokenExpireTime == 0 {
		options.RefreshTokenExpireTime = time.Hour * 24 * 30
	}
	return options
}

type oauth2 struct {
	logger            logger.Logger
	oauth2Persistence storage.OAuth2Persistence
	oauthPersistence  storage.OAuthPersistence
	clientPersistence storage.ClientPersistence
	consentCache      storage.ConsentCache
	authCodeCache     storage.AuthCodeCache
	token             platform.Token
	options           Options
}

func InitOAuth2(logger logger.Logger, oauth2Persistence storage.OAuth2Persistence, oauthPersistence storage.OAuthPersistence, clientPersistence storage.ClientPersistence, consentCache storage.ConsentCache, authCodeCache storage.AuthCodeCache, token platform.Token, options Options) module.OAuth2Module {
	return &oauth2{
		logger:            logger,
		oauth2Persistence: oauth2Persistence,
		oauthPersistence:  oauthPersistence,
		clientPersistence: clientPersistence,
		consentCache:      consentCache,
		authCodeCache:     authCodeCache,
		token:             token,
		options:           options,
	}
}

func (o *oauth2) Authorize(ctx context.Context, authRequestParm dto.AuthorizationRequestParam) (string, errors.AuhtErrResponse, error) {
	if err := authRequestParm.Validate(); err != nil {
		errRsp := errors.AuhtErrResponse{
			Error:            "invalid_request",
			ErrorDescription: strings.TrimSpace(strings.Split(err.Error(), ":")[1]),
		}
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return "", errRsp, err
	}
	client, err := o.oauth2Persistence.GetClient(ctx, authRequestParm.ClientID)
	if err != nil {
		return "", errors.AuhtErrResponse{
			Error:            "invalid_client",
			ErrorDescription: "client not found",
		}, err
	}

	if !o.ContainsRedirectURL(client.RedirectURIs, authRequestParm.RedirectURI) {
		err := errors.ErrInvalidUserInput.New("invalid redirect uri")
		o.logger.Info(ctx, "invalid redirect uri", zap.Error(err))
		return "", errors.AuhtErrResponse{
			Error:            "invalid_redirect_uri",
			ErrorDescription: "invalid redirect uri",
		}, err
	}

	scopes, err := o.oauth2Persistence.GetNamedScopes(ctx, strings.Split(authRequestParm.Scope, " ")...)
	if err != nil || len(scopes) == 0 {
		err := errors.ErrInvalidUserInput.New("invalid scope")
		o.logger.Info(ctx, "invalid scope", zap.Error(err))
		return "", errors.AuhtErrResponse{
			Error:            "invalid_scope",
			ErrorDescription: "invalid scope",
		}, err
	}

	consent := dto.Consent{
		ID:                        uuid.New(),
		AuthorizationRequestParam: authRequestParm,
	}
	if err := o.consentCache.SaveConsent(ctx, consent); err != nil {
		return "", errors.AuhtErrResponse{
			Error:            "server_error",
			ErrorDescription: "failed to save consent",
		}, err
	}

	return consent.ID.String(), errors.AuhtErrResponse{}, nil
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
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		o.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
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

	// check if consent is valid
	consent, err := o.consentCache.GetConsent(ctx, consentId)
	if err != nil || consent.ID.String() != consentId {
		err = errors.ErrNoRecordFound.Wrap(err, "consent not found")
		o.logger.Info(ctx, "consent not found", zap.Error(err), zap.Any("consent-id", consentId))
		return consent, err
	}

	if accessRqResult == "true" {
		var err error
		consent, err = o.consentCache.ChangeStatus(ctx, true, consent)
		if err != nil {
			return dto.Consent{}, err
		}
	} else {
		var err error
		consent, err = o.consentCache.ChangeStatus(ctx, false, consent)
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
	if err := o.authCodeCache.SaveAuthCode(ctx, authCode); err != nil {
		return "", consent.State, err
	}
	return authCode.Code, consent.State, nil
}

func (o *oauth2) Token(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error) {
	if err := param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	grantTypes := map[string]func(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error){
		constant.AuthorizationCode: o.authorizationCodeGrant,
	}

	// Grant processing
	grantHandler := grantTypes[param.GrantType]
	resp, err := grantHandler(ctx, client, param)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *oauth2) authorizationCodeGrant(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error) {
	authcode, err := o.authCodeCache.GetAuthCode(ctx, param.Code)
	if err != nil {
		return nil, err
	}

	user, err := o.oauthPersistence.GetUserByID(ctx, authcode.UserID)
	if err != nil {
		return nil, err
	}

	exists, err := o.oauth2Persistence.AuthHistoryExists(ctx, authcode.Code)
	if err != nil {
		return nil, err
	}
	if exists {

		if err := o.authCodeCache.DeleteAuthCode(ctx, param.Code); err != nil {
			return nil, err
		}

		if err := o.oauth2Persistence.RemoveRefreshToken(ctx, authcode.Code); err != nil {
			return nil, err
		}

		if _, err := o.oauth2Persistence.AddAuthHistory(
			ctx,
			dto.AuthHistory{
				Code:        authcode.Code,
				UserID:      authcode.UserID,
				ClientID:    authcode.ClientID,
				Scope:       authcode.Scope,
				RedirectUri: authcode.RedirectURI,
				Status:      constant.Revoke,
			},
		); err != nil {
			return nil, err
		}

		err := errors.ErrAcessError.New("code already been used")
		o.logger.Info(ctx, "re-use code", zap.Error(err), zap.String("code", authcode.Code))
		return nil, err
	}

	if authcode.ClientID != client.ID {
		err := errors.ErrAuthError.New("client id mismatch")
		o.logger.Warn(ctx, "client id mismatch", zap.Error(err), zap.String("code-client-id", authcode.ClientID.String()), zap.String("given-client-id", client.ID.String()))
		return nil, err
	}

	if param.RedirectURI != "" {
		if authcode.RedirectURI == param.RedirectURI {
			err := errors.ErrAuthError.New("redirect uri mismatch")
			o.logger.Warn(ctx, "redirect uri mismatch", zap.Error(err), zap.String("code-redirect-uri", authcode.RedirectURI), zap.String("given-redirect-uri", param.RedirectURI))
			return nil, err
		}
	}

	accessToken, err := o.token.GenerateAccessToken(ctx, authcode.UserID.String(), o.options.AccessTokenExpireTime)
	if err != nil {
		return nil, err
	}

	idToken, err := o.token.GenerateIdToken(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := o.oauth2Persistence.PersistRefreshToken(ctx, dto.RefreshToken{
		UserID:       authcode.UserID,
		Refreshtoken: o.token.GenerateRefreshToken(ctx),
		ClientID:     authcode.ClientID,
		Scope:        authcode.Scope,
		RedirectUri:  authcode.RedirectURI,
		Code:         authcode.Code,
		ExpiresAt:    time.Now().UTC().Add(o.options.RefreshTokenExpireTime),
	})
	if err != nil {
		return nil, err
	}
	if _, err := o.oauth2Persistence.AddAuthHistory(
		ctx,
		dto.AuthHistory{
			Code:        authcode.Code,
			UserID:      authcode.UserID,
			ClientID:    authcode.ClientID,
			Scope:       authcode.Scope,
			RedirectUri: authcode.RedirectURI,
			Status:      constant.Grant,
		},
	); err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Refreshtoken,
		IDToken:      idToken,
	}, nil
}
