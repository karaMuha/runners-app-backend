package repositories

import (
	"database/sql"
	"net/http"
	"runners/models"
)

type RunnersRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewRunnersRepository(dbHandler *sql.DB) *RunnersRepository {
	return &RunnersRepository{
		dbHandler: dbHandler,
	}
}

func (rr *RunnersRepository) SetTransaction(transaction *sql.Tx) {
	rr.transaction = transaction
}

func (rr *RunnersRepository) GetTransaction() *sql.Tx {
	return rr.transaction
}

func (rr *RunnersRepository) ClearTransaction() {
	rr.transaction = nil
}

func (rr RunnersRepository) QueryCreateRunner(runner *models.Runner) *sql.Row {
	query := `
		INSERT INTO
			runners(first_name, last_name, age, country)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id, is_active`

	row := rr.dbHandler.QueryRow(query, runner.FirstName, runner.LastName, runner.Age, runner.Country)

	return row
}

func (rr RunnersRepository) QueryUpdateRunner(runner *models.Runner) (sql.Result, *models.ResponseError) {
	query := `
		UPDATE
			runners
		SET
			first_name = $1,
			last_name = $2,
			age = $3,
			country = $4
		WHERE
			id = $5`
	res, err := rr.dbHandler.Exec(query, runner.FirstName, runner.LastName, runner.Age, runner.Country, runner.ID)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return res, nil
}

func (rr RunnersRepository) QueryUpdateRunnerResult(runner *models.Runner) (sql.Result, *models.ResponseError) {
	query := `
		UPDATE
			runners
		SET
			personal_best = $1,
			season_best = $2
		WHERE
			id = $3`
	res, err := rr.transaction.Exec(query, runner.PersonalBest, runner.SeasonBest, runner.ID)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return res, nil
}

func (rr RunnersRepository) QueryDeleteRunner(runnerId string) (sql.Result, *models.ResponseError) {
	query := `
		UPDATE
			runners
		SET
			is_active = 'false'
		WHERE
			id = $1`
	res, err := rr.dbHandler.Exec(query, runnerId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return res, nil
}

func (rr RunnersRepository) QueryGetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			runners
		WHERE
			id = $1`
	rows, err := rr.dbHandler.Query(query, runnerId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool
	count := 0

	for rows.Next() {
		count++
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if count == 0 {
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Runner{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		IsActive:     isActive,
		Country:      country,
		PersonalBest: personalBest.String,
		SeasonBest:   seasonBest.String,
	}, nil
}

func (rr RunnersRepository) QueryGetAllRunners() ([]*models.Runner, *models.ResponseError) {
	query := `
		SELECT
			*
		FROM
			runners`
	rows, err := rr.dbHandler.Query(query)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}
		runners = append(runners, runner)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) QueryGetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
	query := `
		SELECT
			id,
			first_name,
			last_name,
			age,
			personal_best,
			season_best
		FROM
			runners
		WHERE
			country = $1
			AND
			is_active = 'true'
		ORDER BY
			personal_best
		LIMIT
			10`
	rows, err := rr.dbHandler.Query(query, country)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName string
	var personalBest, seasonBest sql.NullString
	var age int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     true,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}
		runners = append(runners, runner)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}

func (rr RunnersRepository) QueryGetRunnersByYear(year int) ([]*models.Runner, *models.ResponseError) {
	query := `
		SELECT
			runners.id,
			runners.first_name,
			runners.last_name,
			runners.age,
			runners.is_active,
			runners.country,
			runners.personal_best,
			runners.season_best,
			results.race_result
		FROM
			runners
		INNER JOIN (
			SELECT
				runner_id,
				MIN(race_result) as race_result
			FROM
				results
			WHERE
				year = $1
			GROUP BY
				runner_id
			) results
		ON
			runners.id = results.runner_id
		ORDER BY
			results.race_result
		LIMIT
			10`
	rows, err := rr.dbHandler.Query(query, year)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}

		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}

		runners = append(runners, runner)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return runners, nil
}
