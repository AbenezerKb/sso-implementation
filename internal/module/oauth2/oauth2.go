package oauth2

import (
	"context"
	"net/url"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/state"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/utils"
	"strings"
	"time"

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
}

func InitOAuth2(logger logger.Logger, oauth2Persistence storage.OAuth2Persistence, oauthPersistence storage.OAuthPersistence, clientPersistence storage.ClientPersistence, consentCache storage.ConsentCache, authCodeCache storage.AuthCodeCache, token platform.Token, options Options, scope storage.ScopePersistence) module.OAuth2Module {
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

	scopes, err := o.scopePersistence.GetScopeNameOnly(ctx, strings.Split(authRequestParm.Scope, " ")...)
	if err != nil || scopes == "" {
		err := errors.ErrInvalidUserInput.New("invalid scope")
		o.logger.Info(ctx, "invalid scope", zap.Error(err))
		return "", errors.AuhtErrResponse{
			Error:            "invalid_scope",
			ErrorDescription: "invalid scope",
		}, err
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

func (o *oauth2) GetConsentByID(ctx context.Context, consentID string) (dto.ConsentResponse, error) {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		o.logger.Error(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
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
		Scopes:         requestedscopes,
		Client_Name:    client.Name,
		Client_Logo:    client.LogoURL,
		Client_Type:    client.ClientType,
		Client_Trusted: true,
		Client_ID:      client.ID,
		Approved:       clientStatus,
		User_ID:        user.ID,
	}, nil
}

func (o *oauth2) ApproveConsent(ctx context.Context, consentID string, userID uuid.UUID) (string, error) {
	// check if consent is valid
	consent, err := o.consentCache.GetConsent(ctx, consentID)
	if err != nil {
		err = errors.ErrNoRecordFound.Wrap(err, "consent not found")
		o.logger.Info(ctx, "consent not found", zap.Error(err), zap.Any("consent-id", consentID))
		return "", err
	}

	_, err = o.consentCache.ChangeStatus(ctx, true, consent)
	if err != nil {
		return "", err
	}

	redirectURI, err := url.Parse(consent.RedirectURI)
	if err != nil {
		o.logger.Error(ctx, "invalid redirectURI was found", zap.String("redirect_uri", consent.RedirectURI))
		return "", err
	}

	query := redirectURI.Query()
	query.Set("state", consent.State)

	authCode := dto.AuthCode{
		Code:        utils.GenerateTimeStampedRandomString(25, true),
		Scope:       consent.Scope,
		RedirectURI: consent.RedirectURI,
		ClientID:    consent.ClientID,
		UserID:      userID,
	}
	if err := o.authCodeCache.SaveAuthCode(ctx, authCode); err != nil {
		return "", err
	}

	query.Set("code", authCode.Code)
	if consent.State != "" {
		query.Set("state", consent.State)
	}
	if strings.Contains(consent.ResponseType, "id_token") {
		user, err := o.oauthPersistence.GetUserByID(ctx, userID)
		if err != nil {
			return "", err
		}
		idToken, err := o.token.GenerateIdToken(ctx, user, consent.ClientID.String(), o.options.IDTokenExpireTime)
		if err != nil {
			return "", err
		}
		query.Set("id_token", idToken)
	}

	redirectURI.RawQuery = query.Encode()
	return redirectURI.String(), nil
}

func (o *oauth2) RejectConsent(ctx context.Context, consentID, failureReason string) (string, error) {
	// check if consent is valid
	consent, err := o.consentCache.GetConsent(ctx, consentID)
	if err != nil {
		err = errors.ErrNoRecordFound.Wrap(err, "consent not found")
		o.logger.Info(ctx, "consent not found", zap.Error(err), zap.Any("consent-id", consentID))
		return "", err
	}

	redirectURI, err := url.Parse(consent.RedirectURI)
	if err != nil {
		o.logger.Error(ctx, "invalid redirectURI was found", zap.String("redirect_uri", consent.RedirectURI))
		return "", err
	}

	query := redirectURI.Query()
	query.Set("state", consent.State)
	if failureReason == "" {
		failureReason = "access_denied"
	}
	query.Set("error", failureReason)
	if consent.State != "" {
		query.Set("state", consent.State)
	}

	redirectURI.RawQuery = query.Encode()
	return redirectURI.String(), nil
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

	idToken, err := o.token.GenerateIdToken(ctx, user, authcode.ClientID.String(), o.options.IDTokenExpireTime)
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

func (o oauth2) Logout(ctx context.Context, logoutReqParam dto.LogoutRequest) (string, error) {
	errRedirectUri, err := url.Parse(state.ErrorURL)
	errQuery := errRedirectUri.Query()

	logoutRedirectUri, err := url.Parse(state.LogoutURL)
	logoutQuery := logoutRedirectUri.Query()

	isValid, idToken := o.token.VerifyIdToken(jwt.SigningMethodPS512, logoutReqParam.IDTokenHint)
	if isValid {
		err := errors.ErrInvalidUserInput.New("id_token is invalid")
		o.logger.Error(ctx, "invalid id_token", zap.Error(err), zap.Any("id_token", logoutReqParam.IDTokenHint))
		errQuery.Set("error", "invalid request")
		errQuery.Set("error_description", "no logedin user found")

		errRedirectUri.RawQuery = errQuery.Encode()
		return errRedirectUri.String(), err
	}

	redirectURI, err := url.Parse(logoutReqParam.PostLogoutRedirectUri)
	if err != nil {
		err = errors.ErrInvalidUserInput.New("invalid post logout redirect uri")
		o.logger.Error(ctx, "invalid post logout redirect uri", zap.String("redirect_uri", logoutReqParam.PostLogoutRedirectUri))
		errQuery.Set("error", "invalid post logout redirect uri")
		errQuery.Set("error_description", "post logout redirect uri")

		errRedirectUri.RawQuery = errQuery.Encode()
		return errRedirectUri.String(), err
	}

	// we cleared rf_token
	//
	//
	// o.oauth2Persistence.RemoveRefreshToken(idToken.Subject)

	userID, err := uuid.Parse(idToken.Subject)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		o.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", idToken.Subject))
		errQuery.Set("error", "invalid user input")
		errQuery.Set("error_description", "invalid user input")

		errRedirectUri.RawQuery = errQuery.Encode()
		return errRedirectUri.String(), err
	}

	o.oauthPersistence.RemoveInternalRefreshToken(ctx, userID)
	query := redirectURI.Query()
	query.Set("state", logoutReqParam.State)
	redirectURI.RawQuery = query.Encode()

	logoutQuery.Set("post_logout_uri", redirectURI.String())
	logoutQuery.Set("user_id", idToken.Subject)
	logoutRedirectUri.RawQuery = logoutQuery.Encode()

	return logoutRedirectUri.String(), nil
}
