package persistencedb

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	db2 "sso/internal/constant/model/db"
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

const getUserById = `
SELECT 
id, 
first_name, 
middle_name, 
last_name, 
email, 
phone, 
password, 
user_name, 
gender, 
profile_picture, 
status, 
created_at,
(select v1 from casbin_rule where v0 = cast(users.id as string) limit 1) as role
FROM users WHERE id = $1
`

func (db *PersistenceDB) GetUserByIDWithRole(ctx context.Context, id uuid.UUID) (*dto.User, error) {
	row := db.pool.QueryRow(ctx, getUserById, id)
	var i db2.User
	var role sql.NullString
	err := row.Scan(
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
		&role,
	)
	return &dto.User{
		ID:             i.ID,
		FirstName:      i.FirstName,
		MiddleName:     i.MiddleName,
		LastName:       i.LastName,
		Email:          i.Email.String,
		Phone:          i.Phone,
		Gender:         i.Gender,
		Status:         i.Status.String,
		ProfilePicture: i.ProfilePicture.String,
		CreatedAt:      i.CreatedAt,
		Role:           role.String,
	}, err
}
