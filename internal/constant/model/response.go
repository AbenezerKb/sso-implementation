package model

type Response struct {
	OK       bool           `json:"ok"`
	MetaData *MetaData      `json:"meta_data,omitempty"`
	Data     interface{}    `json:"data,omitempty"`
	Error    *ErrorResponse `json:"error,omitempty"`
}

type MetaData struct {
	PageNo     int         `json:"page_no,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	TotalCount int         `json:"total_count,omitempty"`
	Extra      interface{} `json:"extra,omitempty"`
	Sort       string      `form:"sort" json:"sort,omitempty"`
}

type ErrorResponse struct {
	Code        int      `json:"code"`
	Message     string   `json:"message,omitempty"`
	Description string   `json:"description,omitempty"`
	FieldError  []string `json:"field_error,omitempty"`
}
