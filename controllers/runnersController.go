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

func (rc RunnersController) CreateRunner(ctx *gin.Context) {
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
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		log.Println("Error while umarshling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, responseErr := rc.runnersService.CreateRunner(&runner)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) UpdateRunner(ctx *gin.Context) {
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
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)

	if err != nil {
		metrics.HttpResponsesCounter.WithLabelValues("500").Inc()
		log.Println("Error while umarshling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	responseErr = rc.runnersService.UpdateRunner(&runner)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.Status(http.StatusOK)
}

func (rc RunnersController) DeleteRunner(ctx *gin.Context) {
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

	runnerId := ctx.Param("id")

	responseErr = rc.runnersService.DeleteRunner(runnerId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.Status(http.StatusOK)
}

func (rc RunnersController) GetRunner(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_USER})

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

	runnerId := ctx.Param("id")

	response, responseErr := rc.runnersService.GetRunner(runnerId)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) GetRunnersBatch(ctx *gin.Context) {
	metrics.HttpRequestsCounter.Inc()

	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_USER})

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

	queryParams := ctx.Request.URL.Query()
	country := queryParams.Get("country")
	year := queryParams.Get("year")

	response, responseErr := rc.runnersService.GetRunnersBatch(country, year)

	if responseErr != nil {
		metrics.HttpResponsesCounter.WithLabelValues(strconv.Itoa(responseErr.Status)).Inc()
		ctx.JSON(responseErr.Status, responseErr)
		return
	}

	metrics.HttpResponsesCounter.WithLabelValues("200").Inc()
	ctx.JSON(http.StatusOK, response)
}
