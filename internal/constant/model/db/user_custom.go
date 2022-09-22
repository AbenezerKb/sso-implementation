package db

import (
	"context"
	"sso/platform/utils"
)

func (q *Queries) GetAllUsers(ctx context.Context, pgnFlt string) ([]User, int, error) {
	rows, err := q.db.Query(ctx, utils.ComposeFullFilterSQL(ctx, "users", pgnFlt))
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
