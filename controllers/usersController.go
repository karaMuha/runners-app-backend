package controllers

import (
	"log"
	"net/http"
	"runners/interfaces"
	"runners/metrics"
	"strconv"
)

const ROLE_ADMIN = "admin"
const ROLE_USER = "user"

type UsersController struct {
	usersService interfaces.UsersService
}

func NewUsersController(usersService interfaces.UsersService) *UsersController {
	return &UsersController{
		usersService: usersService,
	}
}

func (uc UsersController) Login(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	username, password, ok := r.BasicAuth()
	if !ok {
		metrics.HttpResponsesCounter.WithLabelValues("400").Inc()
		log.Println("Error while reading credentials")
		http.Error(w, "Error while reading credentials", 400)
		return
	}

	userId, responseErr := uc.usersService.GetUser(username, password)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if userId == "" {
		metrics.HttpResponsesCounter.WithLabelValues("404").Inc()
		http.Error(w, "User not found", 404)
		return
	}

	accessToken, responseErr := uc.usersService.GenerateAccessToken(username, userId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	w.Header().Add("Token", accessToken)
	w.WriteHeader(http.StatusOK)
}

func (uc UsersController) Logout(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsCounter.Inc()

	accessToken := r.Header.Get("Token")

	responseErr := uc.usersService.Logout(accessToken)
	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("204").Inc()
	w.Header().Del("Token")
	w.WriteHeader(http.StatusNoContent)
}
