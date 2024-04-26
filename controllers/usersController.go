package controllers

import (
	"log"
	"net/http"
	"runners/interfaces"
	"runners/metrics"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (uc UsersController) Login(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		metrics.HttpResponsesCounter.WithLabelValues("400").Inc()
		log.Println("Error while readinf credentials")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	accessToken, responseErr := uc.usersService.Login(username, password)
	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, accessToken)
}

func (uc UsersController) Logout(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")

	responseErr := uc.usersService.Logout(accessToken)
	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("204").Inc()
	ctx.Status(http.StatusNoContent)
}
