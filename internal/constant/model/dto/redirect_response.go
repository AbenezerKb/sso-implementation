package dto

type RedirectResponse struct {
	// Location is the url of the client or an error page of the front-end that must be navigated to
	Location string `json:"location"`
}
