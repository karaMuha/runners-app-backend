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

type RunnersController struct {
	runnersService interfaces.RunnersService
	usersService   interfaces.UsersService
}

func NewRunnersController(runnersService interfaces.RunnersService, usersService interfaces.UsersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
		usersService:   usersService,
	}
}

func (rc RunnersController) CreateRunner(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status))
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	var runner models.Runner
	err := json.NewDecoder(r.Body).Decode(&runner)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, "Error while reading request body", http.StatusInternalServerError)
		return
	}

	response, responseErr := rc.runnersService.CreateRunner(&runner)

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

func (rc RunnersController) UpdateRunner(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status))
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	var runner models.Runner
	err := json.NewDecoder(r.Body).Decode(&runner)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, "Error while reading request body", http.StatusInternalServerError)
		return
	}

	rowsAffected, responseErr := rc.runnersService.UpdateRunner(&runner)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if rowsAffected == 0 {
		metrics.HttpResponsesCounter.WithLabelValues("404").Inc()
		http.Error(w, "Runner not found", 404)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	w.WriteHeader(http.StatusOK)
}

func (rc RunnersController) DeleteRunner(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status))
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	runnerId := r.PathValue("id")

	rowsAffected, responseErr := rc.runnersService.DeleteRunner(runnerId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if rowsAffected == 0 {
		metrics.HttpResponsesCounter.WithLabelValues("404").Inc()
		http.Error(w, "Runner not found", 404)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	w.WriteHeader(http.StatusOK)
}

func (rc RunnersController) GetRunner(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN, ROLE_USER})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status))
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	runnerId := r.PathValue("id")

	runner, responseErr := rc.runnersService.GetRunner(runnerId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if runner == nil {
		metrics.HttpResponsesCounter.WithLabelValues("404").Inc()
		http.Error(w, "Runner not found", 404)
		return
	}

	runnersResults, responseErr := rc.runnersService.GetRunnersResults(runner.ID)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	runner.Results = runnersResults

	responseJson, err := json.Marshal(runner)

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

func (rc RunnersController) GetRunnersBatch(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	responseErr := middleware.AuthorizeRequest(r, rc.usersService, []string{ROLE_ADMIN, ROLE_USER})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status))
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	country := r.URL.Query().Get("country")
	year := r.URL.Query().Get("year")

	response, responseErr := rc.runnersService.GetRunnersBatch(country, year)

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
