package interfaces

import "runners/models"

type ResultsServiceInterface interface {
	CreateResult(result *models.Result) (*models.Result, *models.ResponseError)

	UpdateResult(result *models.Result) *models.ResponseError

	DeleteResult(resultId string) *models.ResponseError
}
