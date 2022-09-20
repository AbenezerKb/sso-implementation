package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"sso/internal/constant/model/dto/request_models"
	"sso/platform/logger"
	"strings"
)

func ComposeFilterSQL(ctx context.Context, f request_models.FilterParams, logger logger.Logger) string {
	where := ""
	q := "("
	f.LinkOperator = strings.ToUpper(f.LinkOperator)
	for _, filter := range f.Filter {
		if filter.ColumnField != "*" {
			switch filter.OperatorValue {
			case "=", "!=", ">", ">=", "<", "<=":
				where += fmt.Sprintf("%v %v '%v'", filter.ColumnField, filter.OperatorValue, filter.Value)
			case "is empty":
				where += fmt.Sprintf("%v = 'NULL'", filter.ColumnField)
			case "is not empty":
				where += fmt.Sprintf("%v = 'NOT NULL'", filter.ColumnField)
			case "contains":
				where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, "%"+filter.Value+"%")
			case "equals":
				where += fmt.Sprintf("%v = '%v'", filter.ColumnField, filter.Value)
			case "starts with":
				where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, filter.Value+"%")
			case "ends with":
				where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, "%"+filter.Value)
			case "is":
				where += fmt.Sprintf("%v = '%v'", filter.ColumnField, filter.Value)
			case "is not":
				where += fmt.Sprintf("%v != '%v'", filter.ColumnField, filter.Value)
			case "is after":
				where += fmt.Sprintf("%v > TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			case "is on or after":
				where += fmt.Sprintf("%v >= TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			case "is before":
				where += fmt.Sprintf("%v < TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			case "is on or before":
				where += fmt.Sprintf("%v <= TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			default:
				continue
			}
			where += " " + f.LinkOperator + " "
		} else {
			var searchQ request_models.FilterSearch
			err := json.Unmarshal([]byte(filter.Value), &searchQ)
			if err != nil {
				logger.Info(ctx, "error while parsing a search filter", zap.Error(err), zap.String("filter-value", filter.Value))
				continue
			}
			for _, column := range searchQ.Columns {
				q += fmt.Sprintf("%v ILIKE '%v' OR ", column, "%"+searchQ.Value+"%")
			}
		}
	}

	where = strings.TrimSuffix(where, " "+f.LinkOperator+" ")
	q = strings.TrimSuffix(q, " OR ") + ")"

	sortBy := ""
	for _, v := range f.Sort {
		sortBy += v.Field + " " + v.Sort + ","
	}
	sortBy = strings.TrimSuffix(sortBy, ",")

	return fmt.Sprintf("WHERE %s AND %s ORDER BY %s LIMIT %d OFFSET %d", where, q, sortBy, f.PerPage, f.Page*f.PerPage)
}
