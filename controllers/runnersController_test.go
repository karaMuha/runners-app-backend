package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runners/models"
	"runners/repositories"
	"runners/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func initTestRouter(dbHandler *sql.DB) *http.ServeMux {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	usersService := services.NewUsersService(usersRepository)
	runnersController := NewRunnersController(runnersService, usersService)

	router := http.NewServeMux()
	router.HandleFunc("GET /runner", runnersController.GetRunnersBatch)

	return router
}

func TestGetRunnersResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()

	columns := []string{
		"id",
		"first_name",
		"last_name",
		"age",
		"is_active",
		"country",
		"personal_best",
		"season_best",
	}

	columnsUsers := []string{"user_role"}

	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow(
			"user",
		),
	)

	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).AddRow(
			"5bbff343-29cf-43b5-acc6-6792e3f8074b",
			"Adam",
			"Smith",
			30,
			true,
			"United States",
			"02:00:41",
			"02:13:13",
		).AddRow(
			"6153ccd1-4f64-4d92-8af7-2df6bbdfea66",
			"Sarah",
			"Smith",
			30,
			true,
			"United States",
			"01:18:28",
			"01:18:28",
		),
	)

	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()

	request.Header.Set("Token", "token")
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)

	assert.NotEmpty(t, runners)
	assert.Equal(t, 2, len(runners))
}

func TestGetRunnersErrResponseCountryAndYear(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()

	columnsUsers := []string{"user_role"}

	mock.ExpectQuery("SELECT user_role").WillReturnRows(
		sqlmock.NewRows(columnsUsers).AddRow(
			"user",
		),
	)

	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner?country=france&year=2018", nil)
	recorder := httptest.NewRecorder()

	request.Header.Set("Token", "token")
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
}
