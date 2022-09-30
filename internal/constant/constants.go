package constant

type Context string

const (
	Active      = "ACTIVE"
	Inactive    = "INACTIVE"
	Revoke      = "Revoke"
	Grant       = "Grant"
	User        = "user"
	BearerToken = "Bearer"
	OpenID      = "openid"
	UPDATE      = "UPDATE"
	CREATE      = "CREATE"
	PROMOTE     = "PROMOTE"
)

const (
	AuthorizationCode = "authorization_code"
	RefreshToken      = "refresh_token"
)
