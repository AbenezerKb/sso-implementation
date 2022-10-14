package persistencedb

import (
	"context"
	"database/sql"
	"sso/internal/constant/model/dto"
	"sso/platform/utils"
)

func (db *PersistenceDB) GetAllUsersWithRole(ctx context.Context, pgnFlt string) ([]dto.User, int, error) {
	rows, err := db.pool.Query(ctx, utils.ComposeFullFilterSQL(ctx, "user_role", pgnFlt))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []dto.User
	var totalCount int
	for rows.Next() {
		var i dto.User
		var email, profilePicture, status, role sql.NullString
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.MiddleName,
			&i.LastName,
			&email,
			&i.Phone,
			&i.Gender,
			&profilePicture,
			&status,
			&i.CreatedAt,
			&role,
			&totalCount); err != nil {
			return nil, 0, err
		}
		i.Email = email.String
		i.ProfilePicture = profilePicture.String
		i.Status = status.String
		i.Role = role.String
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
