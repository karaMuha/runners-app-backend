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

	return rs.runnersRepository.QueryCreateRunner(runner)
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError {
	responseErr := validateRunnerId(runner.ID)

	if responseErr != nil {
		return responseErr
	}

	responseErr = validateRunner(runner)

	if responseErr != nil {
		return responseErr
	}

	return rs.runnersRepository.QueryUpdateRunner(runner)
}

func (rs RunnersService) DeleteRunner(runnerId string) *models.ResponseError {
	responseErr := validateRunnerId(runnerId)

	if responseErr != nil {
		return responseErr
	}

	return rs.runnersRepository.QueryDeleteRunner(runnerId)
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
