package state

import (
	"net/url"
)

const (
	ConsentKey  = "consent:%v"
	AuthCodeKey = "authcode:%v"
)

const (
	DefaultPageSize = 10
	LinkOperatorAnd = "AND"
	LinkOperatorOr  = "OR"
)

type URLs struct {
	ErrorURL   *url.URL
	ConsentURL *url.URL
	LogoutURL  *url.URL
}
