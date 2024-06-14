# Runners App
## About this app
This app can be used to store runners and their race results.

The main purpose of this project was to learn how to use golang for backend engineering.

## How run the app
Prerequisites
- Docker installed (I have docker desktop version 4.3 installed on Windows 11)
- Postgres installed (version >=16)

_NOTE: I might refactor this in future to run on docker-compose so it won't be required to have postgres installed, but right now this project does not have a priority_

Clone the code: `git clone https://github.com/karaMuha/runners-app-backend.git`

switch to root dir of the project

Build with docker command
docker build -f Dockerfile . -t runners-app-backend

Run with docker command
docker run -p 8080:8080 runners-app-backend -d

## ToDos
- switch to docker-compose
- provide tests