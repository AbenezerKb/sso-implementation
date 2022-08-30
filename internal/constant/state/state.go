package state

import (
	"net/url"
)

const (
	ConsentKey  = "consent:%v"
	AuthCodeKey = "authcode:%v"
)

type URLs struct {
	ErrorURL   *url.URL
	ConsentURL *url.URL
	LogoutURL  *url.URL
}
