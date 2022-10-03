package dto

type MiniRideResponse struct {
	User   User `json:"user"`
	Exists bool `json:"exists"`
}
