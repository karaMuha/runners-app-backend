package repositories

import (
	"database/sql"
	"net/http"
	"runners/models"
)

type UsersRepository struct {
	dbHandler *sql.DB
}

func NewUsersRepository(dbHandler *sql.DB) *UsersRepository {
	return &UsersRepository{
		dbHandler: dbHandler,
	}
}

func (ur UsersRepository) QueryGetUserId(username string, password string) *sql.Row {
	query := `
				SELECT
					id
				FROM
					users
				WHERE
					username = $1
					AND
					user_password = crypt($2, user_password)`
	row := ur.dbHandler.QueryRow(query, username, password)

	return row
}

func (ur UsersRepository) GetUserRole(accessToken string) (string, *models.ResponseError) {
	query := `
				SELECT
					user_role
				FROM
					users
				WHERE
					access_token = $1`
	rows, err := ur.dbHandler.Query(query, accessToken)

	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var role string
	for rows.Next() {
		err := rows.Scan(&role)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return role, nil
}

func (ur UsersRepository) SetAccessToken(accessToken string, id string) *models.ResponseError {
	query := `
				UPDATE
					users
				SET
					access_token = $1
				WHERE
					id = $2`
	_, err := ur.dbHandler.Exec(query, accessToken, id)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur UsersRepository) RemoveAccessToken(accessToken string) *models.ResponseError {
	query := `
				UPDATE
					users
				SET
					access_token = ''
				WHERE
					access_token = $1`
	_, err := ur.dbHandler.Exec(query, accessToken)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
