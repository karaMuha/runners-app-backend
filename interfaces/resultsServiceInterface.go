package interfaces

import "runners/models"

type ResultsService interface {
	CreateResult(result *models.Result) (*models.Result, *models.ResponseError)

	DeleteResult(resultId string) *models.ResponseError
}
