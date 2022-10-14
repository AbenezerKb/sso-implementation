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
	q := ""
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
			case "startsWith":
				where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, filter.Value+"%")
			case "endsWith":
				where += fmt.Sprintf("%v ILIKE '%v'", filter.ColumnField, "%"+filter.Value)
			case "is":
				where += fmt.Sprintf("%v = '%v'", filter.ColumnField, filter.Value)
			case "not":
				where += fmt.Sprintf("%v != '%v'", filter.ColumnField, filter.Value)
			case "after":
				where += fmt.Sprintf("%v > TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			case "onOrAfter":
				where += fmt.Sprintf("%v >= TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			case "before":
				where += fmt.Sprintf("%v < TIMESTAMPTZ '%v'", filter.ColumnField, filter.Value)
			case "onOrBefore":
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
	q = strings.TrimSuffix(q, " OR ")

	sortBy := ""
	for _, v := range f.Sort {
		sortBy += v.Field + " " + v.Sort + ","
	}
	sortBy = strings.TrimSuffix(sortBy, ",")

	query := ""

	// filters
	if where != "" {
		query += fmt.Sprintf("WHERE %s", where)
	}
	// search filter
	if q != "" {
		if where != "" {
			query += fmt.Sprintf(" AND (%s)", q)
		} else {
			query += fmt.Sprintf("WHERE %s", q)
		}
	}
	// sort
	query += fmt.Sprintf(" ORDER BY %s", sortBy)
	// limit and page number
	// PerPage currently can only be negative or positive integer but not zero
	// if PerPage is negative, no limit and offset it set
	if f.PerPage >= 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", f.PerPage, f.Page*f.PerPage)
	}
	return query
}

func ComposeFullFilterSQL(ctx context.Context, tableName, filterSQL string) string {
	// the COUNT(*) OVER() is used to get the total number of rows without the pagination applied. enjoy!
	return fmt.Sprintf("SELECT *,COUNT(*) OVER() FROM %s %s", tableName, filterSQL)
}
