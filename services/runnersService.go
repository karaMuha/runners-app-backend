package services

import (
	"fmt"
	"net/http"
	"runners/models"
	"runners/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(
	runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunner(runner)

	if responseErr != nil {
		return nil, responseErr
	}

	queryResult := rs.runnersRepository.QueryCreateRunner(runner)

	var runnerId string
	var isActive bool
	err := queryResult.Scan(&runnerId, &isActive)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Runner{
		ID:        runnerId,
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  isActive,
		Country:   runner.Country,
	}, nil
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) (int64, *models.ResponseError) {
	responseErr := validateRunnerId(runner.ID)

	if responseErr != nil {
		return 0, responseErr
	}

	responseErr = validateRunner(runner)

	if responseErr != nil {
		return 0, responseErr
	}

	queryResult, responseErr := rs.runnersRepository.QueryUpdateRunner(runner)

	if responseErr != nil {
		return 0, responseErr
	}

	rowsAffected, err := queryResult.RowsAffected()

	if err != nil {
		return 0, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return rowsAffected, nil
}

func (rs RunnersService) DeleteRunner(runnerId string) (int64, *models.ResponseError) {
	responseErr := validateRunnerId(runnerId)

	if responseErr != nil {
		return 0, responseErr
	}

	queryResult, responseErr := rs.runnersRepository.QueryDeleteRunner(runnerId)

	if responseErr != nil {
		return 0, responseErr
	}

	rowsAffected, err := queryResult.RowsAffected()

	if err != nil {
		return 0, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return rowsAffected, nil
}

func (rs RunnersService) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunnerId(runnerId)

	if responseErr != nil {
		return nil, responseErr
	}

	runner, responseErr := rs.runnersRepository.QueryGetRunner(runnerId)

	if responseErr != nil {
		return nil, responseErr
	}

	results, responseErr := rs.resultsRepository.QueryGetAllRunnersResults(runnerId)

	if responseErr != nil {
		return nil, responseErr
	}

	runner.Results = results

	return runner, nil
}

func (rs RunnersService) GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter can be passed",
			Status:  http.StatusBadRequest,
		}
	}

	if country != "" {
		fmt.Println(country)
		return rs.runnersRepository.QueryGetRunnersByCountry(country)
	}

	if year != "" {
		intYear, err := strconv.Atoi(year)

		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		currentYear := time.Now().Year()

		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}

		return rs.runnersRepository.QueryGetRunnersByYear(intYear)
	}

	return rs.runnersRepository.QueryGetAllRunners()
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if strings.TrimSpace(runner.FirstName) == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}

	if strings.TrimSpace(runner.LastName) == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}

	if runner.Age <= 16 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}

	if strings.TrimSpace(runner.Country) == "" {
		return &models.ResponseError{
			Message: "Invalid country",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}

func validateRunnerId(runnerId string) *models.ResponseError {
	err := uuid.Validate(runnerId)

	if err != nil {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
