package utils

import (
	"context"
	"go.uber.org/zap"
	"sso/internal/constant/model/dto/request_models"
	"sso/platform/logger"
	"testing"
)

type testModel struct {
	data request_models.FilterParams
	want string
}

func TestComposeFilterSQL(t *testing.T) {
	inputs := []testModel{
		{
			data: request_models.FilterParams{
				Sort: []request_models.Sort{
					{
						Field: "name",
						Sort:  "ASC",
					},
				},
				Page:    1,
				PerPage: 20,
				Filter: []request_models.Filter{
					{
						ColumnField:   "name",
						OperatorValue: "is",
						Value:         "kira",
					},
					{
						ColumnField:   "created_at",
						OperatorValue: "is after",
						Value:         "2014-06-06",
					},
					{
						ColumnField: "*",
						Value:       `{"value":"hi","columns":["name", "country", "address"]}`,
					},
				},
				LinkOperator: "AND",
				Total:        0,
				NoSort:       false,
				NoLimit:      false,
			},
			want: "WHERE name = 'kira' AND created_at > TIMESTAMPTZ '2014-06-06' AND (name ILIKE '%hi%' OR country ILIKE '%hi%' OR address ILIKE '%hi%') ORDER BY name ASC LIMIT 20 OFFSET 20",
		},
	}
	log := getLogger(t)

	for _, v := range inputs {
		got := ComposeFilterSQL(context.Background(), v.data, log)
		if got != v.want {
			t.Errorf("\ngot    %s\nwanted %s", got, v.want)
		}
	}
}

func getLogger(t *testing.T) logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	return logger.New(l)

}
