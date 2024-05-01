package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"runners/models"
	"runners/repositories"
	"runners/services"
	"runners/testhelpers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type RunnersControllerTestSuit struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	router      *http.ServeMux
	ctx         context.Context
}

func (suite *RunnersControllerTestSuit) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)

	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	dbHandler, err := sql.Open("postgres", suite.pgContainer.ConnectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = dbHandler.Ping()

	if err != nil {
		log.Fatal(err)
	}

	suite.router = initTestRouter(dbHandler)
}

func (suite *RunnersControllerTestSuit) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	if err != nil {
		log.Fatalf("Error while terminating postgres container: %s", err)
	}
}

func initTestRouter(dbHandler *sql.DB) *http.ServeMux {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	usersService := services.NewUsersService(usersRepository)
	runnersController := NewRunnersController(runnersService, usersService)
	usersController := NewUsersController(usersService)

	router := http.NewServeMux()
	router.HandleFunc("GET /runner", runnersController.GetRunnersBatch)
	router.HandleFunc("POST /login", usersController.Login)

	return router
}

func (suite *RunnersControllerTestSuit) TestGetRunnersResponse() {
	t := suite.T()

	loginRequest, _ := http.NewRequest("POST", "/login", nil)
	loginRequest.SetBasicAuth("user", "user")
	loginRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(loginRecorder, loginRequest)

	require.Equal(t, http.StatusOK, loginRecorder.Result().StatusCode)

	token := loginRecorder.Header().Get("Token")

	request, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()

	request.Header.Set("Token", token)
	suite.router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)

	assert.NotEmpty(t, runners)
	assert.Equal(t, 4, len(runners))
}

func TestRunnersControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RunnersControllerTestSuit))
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
