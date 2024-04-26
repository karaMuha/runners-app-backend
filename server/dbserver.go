package server

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	connectionString := config.GetString("database.connection_string")
	maxIdleConnections := config.GetInt("database.max_idle_connections")
	maxOpenConnections := config.GetInt("database.max_open_connections")
	connectionsMaxLifetime := config.GetDuration("database.connections_max_lifetime")
	driverName := config.GetString("database.driver_name")

	if connectionString == "" {
		log.Fatalf("Database connections tring is missing")
	}

	dbHandler, err := sql.Open(driverName, connectionString)

	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	dbHandler.SetMaxIdleConns(maxIdleConnections)
	dbHandler.SetMaxOpenConns(maxOpenConnections)
	dbHandler.SetConnMaxLifetime(connectionsMaxLifetime)

	err = dbHandler.Ping()

	if err != nil {
		log.Fatalf("Error while validation database: %v", err)
	}

	return dbHandler
}
