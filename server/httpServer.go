package server

import (
	"database/sql"
	"log"
	"net/http"
	"runners/controllers"
	"runners/repositories"
	"runners/services"

	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	server            *http.Server
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
	usersController   *controllers.UsersController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	resultsRepository := repositories.NewResultsRepository(dbHandler)
	usersRepository := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, resultsRepository)
	resultsService := services.NewResultsService(resultsRepository, runnersRepository)
	usersService := services.NewUsersService(usersRepository)
	runnersController := controllers.NewRunnersController(runnersService, usersService)
	resultsController := controllers.NewResultsController(resultsService, usersService)
	usersController := controllers.NewUsersController(usersService)

	router := http.NewServeMux()

	router.HandleFunc("POST /runner", runnersController.CreateRunner)
	router.HandleFunc("PUT /runner", runnersController.UpdateRunner)
	router.HandleFunc("DELETE /runner/{id}", runnersController.DeleteRunner)
	router.HandleFunc("GET /runner/{id}", runnersController.GetRunner)
	router.HandleFunc("GET /runner", runnersController.GetRunnersBatch)

	router.HandleFunc("POST /result", resultsController.CreateResult)
	router.HandleFunc("DELETE /result/{id}", resultsController.DeleteResult)

	router.HandleFunc("POST /login", usersController.Login)
	router.HandleFunc("POST /logout", usersController.Logout)

	server := &http.Server{
		Addr:    config.GetString("http.server_address"),
		Handler: router,
	}

	return HttpServer{
		config:            config,
		server:            server,
		runnersController: runnersController,
		resultsController: resultsController,
		usersController:   usersController,
	}
}

func (hs HttpServer) Start() {
	err := hs.server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
