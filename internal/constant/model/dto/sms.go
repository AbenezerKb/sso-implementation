package dto

type SMS struct {
	To      string `json:"to"`
	Text    string `json:"text"`
	Service string `json:"service"`
}
