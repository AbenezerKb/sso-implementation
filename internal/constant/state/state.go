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
	SortAsc         = "ASC"
	SortDesc        = "DESC"
)

type URLs struct {
	ErrorURL   *url.URL
	ConsentURL *url.URL
	LogoutURL  *url.URL
}

type UploadParams struct {
	FileTypes []FileType
}

type FileType struct {
	Name    string
	Types   []string
	MaxSize int64
}

func (f *FileType) SetValues(values map[string]any) {
	name, ok := values["name"].(string)
	if ok {
		f.Name = name
	}

	types, ok := values["types"].([]any)
	if ok {
		for _, v := range types {
			typeString, ok := v.(string)
			if ok {
				f.Types = append(f.Types, typeString)
			}
		}
	}

	size, ok := values["max_size"].(int)
	if ok {
		f.MaxSize = int64(size)
	}
}
