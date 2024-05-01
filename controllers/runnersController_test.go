package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runners/models"
	"runners/repositories"
	"runners/services"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

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

func TestGetRunnersResponse(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase("runners-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("localtest"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := pgContainer.Terminate(ctx)
		if err != nil {
			t.Fatalf("Failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbHandler, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Fatalf("Failed to connect to database %s", err)
	}

	err = dbHandler.Ping()

	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	router := initTestRouter(dbHandler)
	loginRequest, _ := http.NewRequest("POST", "/login", nil)
	loginRequest.SetBasicAuth("user", "user")
	loginRecorder := httptest.NewRecorder()
	router.ServeHTTP(loginRecorder, loginRequest)

	assert.Equal(t, http.StatusOK, loginRecorder.Result().StatusCode)

	token := loginRecorder.Header().Get("Token")

	request, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()

	request.Header.Set("Token", token)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)

	assert.NotEmpty(t, runners)
	assert.Equal(t, 4, len(runners))
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
