package db

import (
	"context"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
)

func (q *Queries) GetAllUsers(ctx context.Context, pgnFlt db_pgnflt.FilterParams) ([]User, int, error) {
	_, sql := db_pgnflt.GetFilterSQL(pgnFlt)
	rows, err := q.db.Query(ctx, db_pgnflt.GetSelectColumnsQuery([]string{
		"id",
		"first_name",
		"middle_name",
		"last_name",
		"email",
		"phone",
		"password",
		"user_name",
		"gender",
		"profile_picture",
		"status",
		"created_at",
	}, "users", sql))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	var totalCount int
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.MiddleName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.Password,
			&i.UserName,
			&i.Gender,
			&i.ProfilePicture,
			&i.Status,
			&i.CreatedAt,
			&totalCount); err != nil {
			return nil, 0, err
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
