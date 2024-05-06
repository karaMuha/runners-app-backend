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

func (ur UsersRepository) QueryGetUser(username string) *sql.Row {
	query := `
				SELECT
					user_password
				FROM
					users
				WHERE
					username = $1`
	row := ur.dbHandler.QueryRow(query, username)

	return row
}

func (ur UsersRepository) QueryGetUserRole(accessToken string) *sql.Row {
	query := `
				SELECT
					user_role
				FROM
					users
				WHERE
					access_token = $1`
	row := ur.dbHandler.QueryRow(query, accessToken)

	return row
}

func (ur UsersRepository) QuerySetAccessToken(accessToken string, username string) *models.ResponseError {
	query := `
				UPDATE
					users
				SET
					access_token = $1
				WHERE
					username = $2`
	_, err := ur.dbHandler.Exec(query, accessToken, username)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur UsersRepository) QueryRemoveAccessToken(accessToken string) *models.ResponseError {
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
