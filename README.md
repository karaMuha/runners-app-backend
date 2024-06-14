# Runners App
## About this app
This app can be used to store runners and their race results.

The main purpose of this project was to learn how to use golang for backend engineering.

## How run the app
Prerequisites
- Docker installed (I have docker desktop version 4.3 installed on Windows 11)
- Postgres installed (version >=16)

_NOTE: I might refactor this in future to run on docker-compose so it won't be required to have postgres installed and environment variables set up, but right now this project does not have a priority_

- Clone the code: `git clone https://github.com/karaMuha/runners-app-backend.git`

- Open up postgres, create a `runners_db` database and run both sql scripts from directory `dbscripts`

- Check environment variables in `runners.toml` and update if needed

- switch to root dir of the project

- Build with docker command `docker build -f Dockerfile . -t runners-app-backend`

- Run with docker command `docker run -p 8080:8080 runners-app-backend -d`

The entrypoint of the app is `main.go`
On Startup the app will be configured by reading runners.toml in `config.go in package config`. The config will be used to initialize the database `dbserver.go in package server` which will then be used to initialize the http server `httpServer.go in package server`. The http server initializes the logic layers (repositories, services and controllers), sets up the routes and runs the server. 

## Endpoints
_NOTE: if you run the scripts in dbscripts directory you will create an admin (password: admin) and a regular user (password: user) you can use theses users to hit the endpoints_

- POST /login -> Set the credentials (username and password) as basic auth in your request header in order to login
- POST /runner -> Create a runner with following json (Admin route)
```
{
    "first_name": "Max",
    "last_name": "Mustermann",
    "age": 25,
    "country": "Germany"
}
```
- PUT /runner -> Update a runner. Include the the runners ID in the request body (Admin route)
- DELETE /runner/{id} -> Delete runner with corresponding id (Admin route)
- GET /runner/{id} -> Get runner with corresponding id (Admin and User route)
- GET /runner -> Get a batch of runners. Optionally you can specify a year OR a country (Admin and User route)
- POST /result -> Create a race result with following json (Admin route)
```
{
    "runner_id": use id of an existing runner here,
    "race_result": "01:18:10",
    "location": "Germany",
    "position": 6,
    "year": 2024
}
```
- DELETE /result/{id} -> Delete race result with corresponding id (Admin route)
## ToDos
- switch to docker-compose
- provide tests