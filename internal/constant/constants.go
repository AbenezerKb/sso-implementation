package constant

type Context string

const (
	SuperUserRole = "super-user"
)

const (
	Active   = "ACTIVE"
	Inactive = "INACTIVE"
	Pending  = "PENDING"
)

const (
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
	ClientCredentials = "client_credentials"
)

const (
	ClientSecretKey = "the-key-has-to-be-32-bytes-long!"
)

const (
	PromptNone    = "none"
	PromptConsent = "consent"
)
