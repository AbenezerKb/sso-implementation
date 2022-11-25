package persistencedb

import (
	"context"
	"database/sql"
	"fmt"
	"sso/internal/constant/model/dto"
	"strings"
)

// GetUsersByParsedField expects that all values in 'values' are valid for 'fieldName'.
// Use it only if you have made sure the values are valid.
func (db *PersistenceDB) GetUsersByParsedField(ctx context.Context, fieldName string, values []string) ([]dto.User, error) {
	var queries []string // query requests in chunks of 50 users. This is to keep the database load at safe level
	chunks := len(values) / 50
	for i := 0; i < chunks; i += 50 {
		queries = append(queries, strings.Join(values[i:i+50], "','"))
	}
	queries = append(queries, strings.Join(values[chunks*50:], "','")) // the remaining chunk

	var users []dto.User
	for i := 0; i < len(queries); i++ {
		err := func() error { // this is to properly defer rows.Close()
			rows, err := db.pool.Query(ctx,
				fmt.Sprintf(
					`SELECT id,first_name,middle_name,last_name,email,phone,gender,profile_picture,status,created_at
						FROM users WHERE %s in ('%s')`,
					fieldName,
					queries[i]))
			defer rows.Close()

			if err != nil {
				return err
			}

			for rows.Next() {
				var i dto.User
				var email, profilePicture, status sql.NullString
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
				); err != nil {
					return err
				}
				i.Email = email.String
				i.ProfilePicture = profilePicture.String
				i.Status = status.String
				users = append(users, i)
			}
			if err := rows.Err(); err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
