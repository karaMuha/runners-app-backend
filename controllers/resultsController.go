package controllers

import (
	"encoding/json"
	"net/http"
	"runners/interfaces"
	"runners/metrics"
	"runners/middleware"
	"runners/models"
	"strconv"
)

type ResultsController struct {
	resultsService interfaces.ResultsServiceInterface
	usersService   interfaces.UsersService
}

func NewResultsController(resultsService interfaces.ResultsServiceInterface, usersService interfaces.UsersService) *ResultsController {
	return &ResultsController{
		resultsService: resultsService,
		usersService:   usersService,
	}
}

func (rc ResultsController) CreateResult(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	var result models.Result
	err := json.NewDecoder(r.Body).Decode(&result)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, responseErr := rc.resultsService.CreateResult(&result)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	responseJson, err := json.Marshal(response)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}

func (rc ResultsController) DeleteResult(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	resultId := r.PathValue("id")
	responseErr = rc.resultsService.DeleteResult(resultId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("204").Inc()
	w.WriteHeader(http.StatusNoContent)
}

func (rc ResultsController) UpdateResult(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	var result models.Result
	err := json.NewDecoder(r.Body).Decode(&result)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseErr = rc.resultsService.UpdateResult(&result)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	w.WriteHeader(http.StatusOK)
}
