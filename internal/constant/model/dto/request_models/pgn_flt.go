package request_models

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/url"
	"sso/internal/constant/state"
	"sso/platform/utils/collection"
	"sso/platform/utils/reflect"
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
	Sort         string `json:"sort" form:"sort"`
	Filter       string `json:"filter" form:"filter"`
	Page         int    `json:"page" form:"page"`
	PerPage      int    `json:"per_page" form:"per_page"`
	LinkOperator string `json:"link_operator" form:"link_operator"`
}

// ToFilterParams returns the FilterParam object this pgnFltQueryParams holds.
// model is the type of database model this filter is being applied to
func (q *PgnFltQueryParams) ToFilterParams(model interface{}) (FilterParams, error) {
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
		for k := range res.Sort {
			res.Sort[k].Sort = strings.ToUpper(res.Sort[k].Sort)
			if err := validation.Validate(res.Sort[k].Sort, validation.In(state.SortAsc, state.SortDesc)); err != nil {
				res.Sort[k].Sort = state.SortDesc
			}
		}
	} else {
		res.Sort = []Sort{
			{
				Field: "created_at",
				Sort:  state.SortDesc,
			},
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
		for _, filter := range filters {
			if filter.ColumnField != "*" {
				if !collection.Contains[string](filter.ColumnField, reflect.GetJSONFieldNames(model)) {
					return FilterParams{}, fmt.Errorf("invalid filter column %s", filter.ColumnField)
				}
			}
		}
	}

	if q.PerPage == 0 {
		res.PerPage = state.DefaultPageSize
	} else {
		res.PerPage = q.PerPage
	}
	q.LinkOperator = strings.ToUpper(q.LinkOperator)
	if err := validation.Validate(q.LinkOperator, validation.Required, validation.In(state.LinkOperatorAnd, state.LinkOperatorOr)); err == nil {
		res.LinkOperator = q.LinkOperator
	} else {
		res.LinkOperator = state.LinkOperatorOr
	}

	res.Filter = filters
	return res, nil
}
