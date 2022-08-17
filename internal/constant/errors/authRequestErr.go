package errors

type AuhtErrResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
