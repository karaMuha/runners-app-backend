package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners/interfaces"
	"runners/metrics"
	"runners/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResultsController struct {
	resultsService interfaces.ResultsService
	usersService   interfaces.UsersService
}

func NewResultsController(resultsService interfaces.ResultsService, usersService interfaces.UsersService) *ResultsController {
	return &ResultsController{
		resultsService: resultsService,
		usersService:   usersService,
	}
}

func (rc ResultsController) CreateResult(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		metrics.HttpResponsesCounter.WithLabelValues("401").Inc()
		ctx.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var result models.Result
	err = json.Unmarshal(body, &result)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		log.Println("Error while unmarshling creates result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.resultsService.CreateResult(&result)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc ResultsController) DeleteResult(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	if !auth {
		metrics.HttpResponsesCounter.WithLabelValues("401").Inc()
		ctx.Status(http.StatusUnauthorized)
		return
	}

	resultId := ctx.Param("id")
	responseErr = rc.resultsService.DeleteResult(resultId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("204").Inc()
	ctx.Status(http.StatusNoContent)
}
