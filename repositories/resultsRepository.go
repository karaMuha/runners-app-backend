package repositories

import (
	"database/sql"
	"net/http"
	"runners/models"
)

type ResultsRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewResultsRepository(dbHandler *sql.DB) *ResultsRepository {
	return &ResultsRepository{
		dbHandler: dbHandler,
	}
}

func (rr *ResultsRepository) SetTransaction(transaction *sql.Tx) {
	rr.transaction = transaction
}

func (rr *ResultsRepository) GetTransaction() *sql.Tx {
	return rr.transaction
}

func (rr *ResultsRepository) ClearTransaction() {
	rr.transaction = nil
}

func (rr ResultsRepository) QueryCreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	query := `
		INSERT INTO
			results(runner_id, race_result, location, position, year)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id`
	rows, err := rr.dbHandler.Query(query, result.RunnerID, result.RaceResult, result.Location, result.Position, result.Year)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var resultId string
	for rows.Next() {
		err := rows.Scan(&resultId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         resultId,
		RunnerID:   result.RunnerID,
		RaceResult: result.RaceResult,
		Location:   result.Location,
		Position:   result.Position,
		Year:       result.Year,
	}, nil
}

func (rr ResultsRepository) QueryUpdateResult(result *models.Result) *models.ResponseError {
	query := `
		UPDATE
			results
		SET
			race_result = $1,
			location = $2,
			position = $3,
			year = $4
		WHERE
			result_id = $5
	`
	res, err := rr.dbHandler.Exec(query, result.RaceResult, result.Location, result.Position, result.Year, result.ID)

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Result not found",
			Status:  http.StatusNotFound,
		}
	}

	return nil
}

func (rr ResultsRepository) QueryDeleteResult(resultId string) (*models.Result, *models.ResponseError) {
	query := `
		DELETE FROM
			results
		WHERE
			id = $1
		RETURNING
			runner_id, race_result, year`
	rows, err := rr.dbHandler.Query(query, resultId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var runnerId, raceResult string
	var year int
	for rows.Next() {
		err := rows.Scan(&runnerId, &raceResult, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Result{
		ID:         resultId,
		RunnerID:   runnerId,
		RaceResult: raceResult,
		Year:       year,
	}, nil
}

func (rr ResultsRepository) QueryGetAllRunnersResults(runnerId string) ([]*models.Result, *models.ResponseError) {
	query := `
		SELECT
			id, race_result, location, position, year
		FROM
			results
		WHERE 
			runner_id = $1`
	rows, err := rr.dbHandler.Query(query, runnerId)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	results := make([]*models.Result, 0)
	var id, raceResult, location string
	var position, year int

	for rows.Next() {
		err := rows.Scan(&id, &raceResult, &location, &position, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		result := &models.Result{
			ID:         id,
			RunnerID:   runnerId,
			RaceResult: raceResult,
			Location:   location,
			Position:   position,
			Year:       year,
		}
		results = append(results, result)
	}

	err = rows.Err()
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return results, nil
}

func (rr ResultsRepository) QueryGetPersonalBestResults(runnerId string) (string, *models.ResponseError) {
	query := `
		SELECT
			MIN(race_result)
		FROM
			results
		WHERE
			runner_id = $1`
	rows, err := rr.dbHandler.Query(query, runnerId)

	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var raceResult string
	for rows.Next() {
		err := rows.Scan(&raceResult)
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

	return raceResult, nil
}

func (rr ResultsRepository) QueryGetSeasonBestResults(runnerId string, year int) (string, *models.ResponseError) {
	query := `
		SELECT
			MIN(race_result)
		FROM
			results
		WHERE
			runner_id = $1
			AND
			year = $2`
	rows, err := rr.dbHandler.Query(query, runnerId, year)

	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var raceResult string
	for rows.Next() {
		err := rows.Scan(&raceResult)
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

	return raceResult, nil
}
