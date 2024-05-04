package interfaces

import "runners/models"

type RunnersService interface {
	CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError)

	UpdateRunner(runner *models.Runner) (int64, *models.ResponseError)

	DeleteRunner(runnerId string) (int64, *models.ResponseError)

	GetRunner(runnerId string) (*models.Runner, *models.ResponseError)

	GetRunnersBatch(country string, year string) ([]*models.Runner, *models.ResponseError)
}
