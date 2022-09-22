package request_models

import (
	"reflect"
	"testing"
)

type User struct {
	Name string
}

func TestPgnFltQueryParams_ToFilterParams(t *testing.T) {
	query := &PgnFltQueryParams{
		Sort:   `[{"field":"name", "sort":"desc"}]`,
		Filter: `[{"column_field":"name","operator_value":"contains","value":"bek"}, {"column_field":"*","value":"b"}]`,

		Page:         0,
		PerPage:      10,
		LinkOperator: "AND",
	}

	want := FilterParams{
		Sort: []Sort{
			{
				Field: "name",
				Sort:  "DESC",
			},
		},
		Page:    0,
		PerPage: 10,
		Filter: []Filter{
			{
				ColumnField:   "name",
				OperatorValue: "contains",
				Value:         "bek",
			},
			{
				ColumnField: "*",
				Value:       "b",
			},
		},
		LinkOperator: "AND",
	}
	got, err := query.ToFilterParams(User{})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
