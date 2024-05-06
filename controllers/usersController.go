package controllers

import (
	"net/http"
	"runners/interfaces"
	"runners/metrics"
	"strconv"

	"golang.org/x/crypto/bcrypt"
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
		http.Error(w, "Error while reading credentials", 400)
		return
	}

	userPassword, responseErr := uc.usersService.GetUser(username)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		http.Error(w, responseErr.Message, responseErr.Status)
		return
	}

	if userPassword == "" {
		metrics.HttpResponsesCounter.WithLabelValues("404").Inc()
		http.Error(w, "User not found", 404)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("401").Inc()
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	accessToken, responseErr := uc.usersService.GenerateAccessToken(username)

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
