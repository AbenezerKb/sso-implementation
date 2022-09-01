package request_models

type RevokeClientBody struct {
	// ClientID is the id of the client to be revoked access of.
	// It is a required field
	ClientID string `json:"client_id"`
}
