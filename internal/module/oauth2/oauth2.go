package oauth2

import (
	"context"
	"fmt"
	"net/url"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/constant/state"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/utils"
	"strings"
	"time"

	"github.com/joomcode/errorx"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Options struct {
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
	IDTokenExpireTime      time.Duration
}

func SetOptions(options Options) Options {
	if options.AccessTokenExpireTime == 0 {
		options.AccessTokenExpireTime = time.Minute * 10
	}
	if options.RefreshTokenExpireTime == 0 {
		options.RefreshTokenExpireTime = time.Hour * 24 * 30
	}
	if options.IDTokenExpireTime == 0 {
		options.IDTokenExpireTime = time.Minute * 10
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
	scopePersistence  storage.ScopePersistence
	urls              state.URLs
}

func InitOAuth2(logger logger.Logger, oauth2Persistence storage.OAuth2Persistence, oauthPersistence storage.OAuthPersistence, clientPersistence storage.ClientPersistence, consentCache storage.ConsentCache, authCodeCache storage.AuthCodeCache, token platform.Token, options Options, scope storage.ScopePersistence, urls state.URLs) module.OAuth2Module {
	return &oauth2{
		logger:            logger,
		oauth2Persistence: oauth2Persistence,
		oauthPersistence:  oauthPersistence,
		clientPersistence: clientPersistence,
		consentCache:      consentCache,
		authCodeCache:     authCodeCache,
		token:             token,
		options:           options,
		scopePersistence:  scope,
		urls:              urls,
	}
}

func (o *oauth2) Authorize(ctx context.Context, authRequestParm dto.AuthorizationRequestParam, requestOrigin string, bindError *errorx.Error) string {
	if bindError != nil {
		o.logger.Info(ctx, "error while binding to query", zap.Error(bindError))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             bindError.Message(),
			"error_description": bindError.Error(),
		})
	}

	if er := authRequestParm.Validate(); er != nil {
		err := errors.ErrInvalidUserInput.Wrap(er, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             "invalid_request",
			"error_description": strings.TrimSpace(strings.Split(er.Error(), ":")[1]),
		})
	}
	redirectURI, err := url.Parse(authRequestParm.RedirectURI)
	if err != nil {
		o.logger.Info(ctx, "error parsing redirect uri", zap.Error(err))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             "invalid_redirect_uri",
			"error_description": err.Error(),
		})
	}

	client, err := o.clientPersistence.GetClientByID(ctx, authRequestParm.ClientID)
	if err != nil {
		return utils.GenerateRedirectString(redirectURI, map[string]string{
			"error":             "invalid_client",
			"error_description": "client not found",
		})
	}

	if !o.ContainsRedirectURL(client.RedirectURIs, authRequestParm.RedirectURI) {
		err := errors.ErrInvalidUserInput.New("invalid redirect uri")
		o.logger.Info(ctx, "invalid redirect uri", zap.Error(err))

		return utils.GenerateRedirectString(redirectURI, map[string]string{
			"error":             "invalid_redirect_uri",
			"error_description": "invalid redirect uri",
		})
	}

	scopes, err := o.scopePersistence.GetScopeNameOnly(ctx, strings.Split(authRequestParm.Scope, " ")...)
	if err != nil || scopes == "" {
		err := errors.ErrInvalidUserInput.New("invalid scope")
		o.logger.Info(ctx, "invalid scope", zap.Error(err))

		return utils.GenerateRedirectString(redirectURI, map[string]string{
			"error":             "invalid_scope",
			"error_description": "invalid scope",
		})
	}

	consent := dto.Consent{
		ID: uuid.New(),
		AuthorizationRequestParam: dto.AuthorizationRequestParam{
			ClientID:     client.ID,
			Scope:        scopes,
			RedirectURI:  authRequestParm.RedirectURI,
			State:        authRequestParm.State,
			ResponseType: authRequestParm.ResponseType,
			Prompt:       authRequestParm.Prompt,
		},
		RequestOrigin: requestOrigin,
	}
	if err := o.consentCache.SaveConsent(ctx, consent); err != nil {
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             "server_error",
			"error_description": "failed to save consent",
		})
	}
	prompt := "consent"
	if authRequestParm.Prompt != "" {
		prompt = authRequestParm.Prompt
	}

	return utils.GenerateRedirectString(o.urls.ConsentURL, map[string]string{
		"consentId": consent.ID.String(),
		"prompt":    prompt,
	})
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

func (o *oauth2) GetConsentByID(ctx context.Context, consentID string) (dto.ConsentResponse, error) {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		o.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return dto.ConsentResponse{}, err
	}

	consent, err := o.consentCache.GetConsent(ctx, consentID)
	if err != nil {
		err = errors.ErrNoRecordFound.Wrap(err, "consent not found")
		o.logger.Info(ctx, "consent not found", zap.Error(err))
		return dto.ConsentResponse{}, err
	}

	// get client
	client, err := o.clientPersistence.GetClientByID(ctx, consent.ClientID)
	if err != nil {
		return dto.ConsentResponse{}, err
	}

	// get scopes
	requestedscopes, err := o.scopePersistence.GetListedScopes(ctx, strings.Split(consent.Scope, " ")...)
	if err != nil {
		return dto.ConsentResponse{}, err
	}

	// get user
	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		o.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return dto.ConsentResponse{}, err
	}
	user, err := o.oauthPersistence.GetUserByID(ctx, userID)
	if err != nil {
		return dto.ConsentResponse{}, err
	}

	clientStatus := true
	check, refreshToken, err := o.oauth2Persistence.CheckIfUserGrantedClient(ctx, userID, client.ID)
	if err != nil {
		return dto.ConsentResponse{}, err
	}

	grantedScopes := utils.StringToArray(refreshToken.Scope)
	if check {
		for _, rs := range requestedscopes {
			if !utils.ContainsValue(rs.Name, grantedScopes) {
				clientStatus = false
				break
			}
		}
	} else {
		clientStatus = false
	}
	return dto.ConsentResponse{
		Scopes:        requestedscopes,
		ClientName:    client.Name,
		ClientLogo:    client.LogoURL,
		ClientType:    client.ClientType,
		ClientTrusted: false,
		ClientID:      client.ID,
		Approved:      clientStatus,
		UserID:        user.ID,
	}, nil
}

func (o *oauth2) ApproveConsent(ctx context.Context, consentID string, userID uuid.UUID, opbs string, bindError *errorx.Error) string {
	if bindError != nil {
		o.logger.Info(ctx, "error while binding to query", zap.Error(bindError))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error": bindError.Message(),
		})
	}
	// check if consent is valid
	consent, err := o.consentCache.GetConsent(ctx, consentID)
	if err != nil {
		o.logger.Info(ctx, "consent not found", zap.Error(err), zap.Any("consent-id", consentID))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":       "consent not found",
			"description": err.Error(),
		})
	}

	redirectURI, err := url.Parse(consent.RedirectURI)
	if err != nil {
		o.logger.Error(ctx, "invalid redirectURI was found", zap.Error(err), zap.String("redirect_uri", consent.RedirectURI))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error": "invalid redirectURI was found",
		})
	}

	authCode := dto.AuthCode{
		Code:        utils.GenerateTimeStampedRandomString(25, false),
		Scope:       consent.Scope,
		RedirectURI: consent.RedirectURI,
		ClientID:    consent.ClientID,
		UserID:      userID,
		State:       consent.State,
	}
	if err := o.authCodeCache.SaveAuthCode(ctx, authCode); err != nil {
		errx := errorx.Cast(err)
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":       errx.Message(),
			"description": errx.Error(),
		})
	}

	queries := map[string]string{}
	queries["code"] = authCode.Code
	if consent.State != "" {
		queries["state"] = consent.State
	}

	// calculate session state
	sessionState := utils.CalculateSessionState(authCode.ClientID.String(), consent.RequestOrigin, opbs, utils.GenerateRandomString(20, true))
	queries["session_state"] = sessionState

	return utils.GenerateRedirectString(redirectURI, queries)
}

func (o *oauth2) RejectConsent(ctx context.Context, consentID, failureReason string, bindError *errorx.Error) string {
	if bindError != nil {
		o.logger.Info(ctx, "error while binding to query", zap.Error(bindError))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error": bindError.Message(),
		})
	}
	// check if consent is valid
	consent, err := o.consentCache.GetConsent(ctx, consentID)
	if err != nil {
		o.logger.Info(ctx, "consent not found", zap.Error(err), zap.Any("consent-id", consentID))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":       "consent not found",
			"description": err.Error(),
		})
	}

	redirectURI, err := url.Parse(consent.RedirectURI)
	if err != nil {
		o.logger.Error(ctx, "invalid redirectURI was found", zap.Error(err), zap.String("redirect_uri", consent.RedirectURI))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error": "invalid redirectURI was found",
		})
	}

	queries := map[string]string{}
	if failureReason == "" {
		failureReason = "unknown error"
	}
	queries["error"] = failureReason
	if consent.State != "" {
		queries["state"] = consent.State
	}

	return utils.GenerateRedirectString(redirectURI, queries)
}

func (o *oauth2) Token(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error) {
	if err := param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	grantTypes := map[string]func(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error){
		constant.AuthorizationCode: o.authorizationCodeGrant,
		constant.RefreshToken:      o.refreshToken,
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

	refreshToken, err := o.oauth2Persistence.PersistRefreshToken(ctx, dto.RefreshToken{
		UserID:       authcode.UserID,
		RefreshToken: o.token.GenerateRefreshToken(ctx),
		ClientID:     authcode.ClientID,
		Scope:        authcode.Scope,
		RedirectUri:  authcode.RedirectURI,
		Code:         authcode.Code,
		ExpiresAt:    time.Now().Add(o.options.RefreshTokenExpireTime),
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
	tokenResponse := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshToken,
		TokenType:    constant.BearerToken,
		ExpiresIn:    fmt.Sprintf("%vs", o.options.AccessTokenExpireTime.Seconds()),
	}
	if utils.ContainsValue(constant.OpenID, utils.StringToArray(authcode.Scope)) {

		user, err := o.oauthPersistence.GetUserByID(ctx, authcode.UserID)
		if err != nil {
			return nil, err
		}

		idToken, err := o.token.GenerateIdToken(ctx, user, client.ID.String(), o.options.IDTokenExpireTime)
		if err != nil {
			return nil, err
		}
		tokenResponse.IDToken = idToken

	}
	return tokenResponse, nil
}

func (o *oauth2) refreshToken(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error) {
	oldRefreshToken, err := o.oauth2Persistence.GetRefreshToken(ctx, param.RefreshToken)
	if err != nil {
		return nil, err
	}
	if oldRefreshToken.ClientID != client.ID {
		err := errors.ErrAuthError.New("client id mismatch")
		o.logger.Warn(ctx, "client id mismatch", zap.Error(err), zap.String("refresh-token-client-id", oldRefreshToken.ClientID.String()), zap.String("given-client-id", client.ID.String()))
		return nil, err
	}

	if time.Now().After(oldRefreshToken.ExpiresAt) {
		if err := o.oauth2Persistence.RemoveRefreshToken(ctx, oldRefreshToken.RefreshToken); err != nil {
			return nil, err
		}

		err := errors.ErrAuthError.New("refresh token expired")
		o.logger.Warn(ctx, "token expired", zap.Error(err), zap.String("refresh token", oldRefreshToken.RefreshToken))
		return nil, err
	}

	accessToken, err := o.token.GenerateAccessTokenForClient(ctx, oldRefreshToken.UserID.String(), oldRefreshToken.ClientID.String(), oldRefreshToken.Scope, o.options.AccessTokenExpireTime)
	if err != nil {
		return nil, err
	}

	if err := o.oauth2Persistence.RemoveRefreshToken(ctx, oldRefreshToken.RefreshToken); err != nil {
		return nil, err
	}

	newRefreshToken, err := o.oauth2Persistence.PersistRefreshToken(ctx, dto.RefreshToken{
		UserID:       oldRefreshToken.UserID,
		RefreshToken: o.token.GenerateRefreshToken(ctx),
		ClientID:     oldRefreshToken.ClientID,
		Scope:        oldRefreshToken.Scope,
		RedirectUri:  oldRefreshToken.RedirectUri,
		ExpiresAt:    time.Now().Add(o.options.RefreshTokenExpireTime),
	})
	if err != nil {
		return nil, err
	}
	tokenResponse := &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.RefreshToken,
		TokenType:    constant.BearerToken,
		ExpiresIn:    fmt.Sprintf("%vs", o.options.AccessTokenExpireTime.Seconds()),
	}

	if utils.ContainsValue(constant.OpenID, utils.StringToArray(newRefreshToken.Scope)) {

		user, err := o.oauthPersistence.GetUserByID(ctx, newRefreshToken.UserID)
		if err != nil {
			return nil, err
		}

		idToken, err := o.token.GenerateIdToken(ctx, user, client.ID.String(), o.options.IDTokenExpireTime)
		if err != nil {
			return nil, err
		}
		tokenResponse.IDToken = idToken

	}
	return tokenResponse, nil
}

func (o *oauth2) Logout(ctx context.Context, logoutReqParam dto.LogoutRequest, bindError *errorx.Error) string {
	if bindError != nil {
		o.logger.Info(ctx, "error while binding to query", zap.Error(bindError))
		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error": bindError.Message(),
		})
	}

	if err := logoutReqParam.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid request", zap.Error(err))

		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             "invalid request",
			"error_description": "invalid request",
		})
	}

	isValid, idToken := o.token.VerifyIdToken(jwt.SigningMethodPS512, logoutReqParam.IDTokenHint)
	if !isValid {
		err := errors.ErrInvalidUserInput.New("id_token is invalid")
		o.logger.Info(ctx, "invalid id_token", zap.Error(err), zap.Any("id_token", logoutReqParam.IDTokenHint))

		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             "invalid request",
			"error_description": "no logedin user found",
		})
	}

	postLogoutgredirectURI, err := url.Parse(logoutReqParam.PostLogoutRedirectUri)
	if err != nil {
		err = errors.ErrInvalidUserInput.New("invalid post logout redirect uri")
		o.logger.Info(ctx, "invalid post logout redirect uri", zap.String("redirect_uri", logoutReqParam.PostLogoutRedirectUri))

		return utils.GenerateRedirectString(o.urls.ErrorURL, map[string]string{
			"error":             "invalid post logout redirect uri",
			"error_description": "post logout redirect uri is invalid",
		})
	}

	return utils.GenerateRedirectString(o.urls.LogoutURL, map[string]string{
		"post_logout_redirect_uri": postLogoutgredirectURI.String(),
		"state":                    logoutReqParam.State,
		"user_id":                  idToken.Subject,
	})

}

func (o *oauth2) RevokeClient(ctx context.Context, clientBody request_models.RevokeClientBody) error {
	if err := clientBody.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid user input", zap.Error(err))
		return err
	}
	userIDString, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		o.logger.Error(ctx, "expected to find x-user-id on context", zap.Error(err), zap.Any("user-id", userIDString))
		return err
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		o.logger.Error(ctx, "unexpected parse error for user id in context", zap.Error(err), zap.String("user-id", userIDString))
		return err
	}
	clientID, err := uuid.Parse(clientBody.ClientID)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid client_id")
		o.logger.Info(ctx, "invalid client_id was provided to revoke client access", zap.Error(err), zap.Any("user-id", userID), zap.String("client-id", clientBody.ClientID))
		return err
	}
	// check refresh token with client id and user id
	refreshToken, err := o.oauth2Persistence.GetRefreshTokenOfClientByUserID(ctx, userID, clientID)
	if err != nil {
		if errorx.IsOfType(err, errors.ErrNoRecordFound) {
			err := errors.ErrInvalidUserInput.Wrap(err, "no client access found")
			return err
		}
		return err
	}

	// delete the refresh token
	err = o.oauth2Persistence.RemoveRefreshToken(ctx, refreshToken.RefreshToken)
	if err != nil {
		return err
	}

	// create an auth history
	_, err = o.oauth2Persistence.AddAuthHistory(ctx, dto.AuthHistory{
		Code:        refreshToken.Code,
		UserID:      userID,
		ClientID:    clientID,
		Scope:       refreshToken.Scope,
		Status:      constant.Revoke,
		RedirectUri: refreshToken.RedirectUri,
	})
	if err != nil {
		return err
	}

	return nil
}
