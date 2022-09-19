package request_models

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/url"
	"sso/internal/constant/state"
	"strings"
)

type Sort struct {
	Field string `json:"field"`
	Sort  string `json:"sort"`
}

type FilterParams struct {
	Sort         []Sort   `json:"sort"`
	Page         int      `json:"page"`
	PerPage      int      `json:"per_page"`
	Filter       []Filter `json:"filter"`
	LinkOperator string   `json:"link_operator"`
	Total        int      `json:"total"`
	NoSort       bool     `json:"-"`
	NoLimit      bool     `json:"-"`
}

type Filter struct {
	ColumnField   string `json:"column_field"`
	OperatorValue string `json:"operator_value"`
	Value         string `json:"value"`
}

type FilterSearch struct {
	Columns []string `json:"columns"`
	Value   string   `json:"value"`
}
type PgnFltQueryParams struct {
	Sort         string `json:"sort"`
	Filter       string `json:"filter"`
	Page         int    `json:"page"`
	PerPage      int    `json:"per_page"`
	LinkOperator string `json:"linkOperator"`
}

// Get returns the FilterParam object this pgnFltQueryParams holds
func (q *PgnFltQueryParams) toFilterParams() (FilterParams, error) {
	res := FilterParams{}
	if q.Sort != "" {
		sortString, err := url.QueryUnescape(q.Sort)
		if err != nil {
			return FilterParams{}, err
		}

		err = json.Unmarshal([]byte(sortString), &res.Sort)
		if err != nil {
			return FilterParams{}, err
		}
	}
	if q.Page < 0 {
		res.Page = 0
	} else {
		res.Page = q.Page
	}

	var filters []Filter
	if q.Filter != "" {
		err := json.Unmarshal([]byte(q.Filter), &filters)
		if err != nil {
			return FilterParams{}, err
		}
	}

	if q.PerPage <= 0 {
		res.PerPage = state.DefaultPageSize
	} else {
		res.PerPage = q.PerPage
	}
	q.LinkOperator = strings.ToUpper(q.LinkOperator)
	if err := validation.Validate(q.LinkOperator, validation.In(state.LinkOperatorAnd, state.LinkOperatorOr)); err == nil {
		res.LinkOperator = q.LinkOperator
	} else {
		res.LinkOperator = state.LinkOperatorOr
	}

	res.Filter = filters
	return res, nil
}
