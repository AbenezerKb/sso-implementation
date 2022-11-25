package dto

// RSAPIUsersResponse contains response data for get users by id or phone
type RSAPIUsersResponse struct {
	// IDs is the users fetched using ids
	IDs []User `json:"ids,omitempty"`
	// Phones is the users fetched using phones
	Phones []User `json:"phones,omitempty"`
}
