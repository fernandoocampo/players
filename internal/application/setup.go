package application

import (
	"log/slog"
	"os"
	"strconv"
)

// Settings contains the parameters required for this service to work.
type settings struct {
	// database parameters.
	repository databaseSetup
	// LogLevel it could be 'production' or 'development'.
	logLevel string
	// tracerServiceURL contains the OTL service URL.
	tracerServiceURL string
	// password generation cost
	passwordGenerationCost int
	// port for web server
	webServerPort int
	// port for grpc server
	grpcServerPort int
	// timeout value used to wait for notifier to push events in event bus.
	timeoutToPublishSec int
}

type databaseSetup struct {
	dbName   string
	host     string
	player   string
	password string
	port     int
}

const (
	logLevelEnvVar               = "PLAYERS_LOG_LEVEL"
	webServerPortEnvVar          = "PLAYERS_WEB_SERVER_PORT"
	grpcServerPortEnvVar         = "PLAYERS_GRPC_SERVER_PORT"
	passwordGenerationCostEnvVar = "PLAYERS_PASSWORD_GENERATION_COST"
	postgresDBNameEnvVar         = "PLAYERS_POSTGRES_DB"
	postgresHostEnvVar           = "PLAYERS_POSTGRES_HOST"
	postgresPlayerEnvVar         = "PLAYERS_POSTGRES_PLAYER"
	postgresPasswordEnvVar       = "PLAYERS_POSTGRES_PASSWORD"
	postgresPortEnvVar           = "PLAYERS_POSTGRES_PORT"
	timeoutToPublishSecEnvVar    = "PLAYERS_TIMEOUT_TO_PUBLISH_SEC"
	tracerServiceURL             = "PLAYERS_TRACER_SERVICE_URL"
)

// log levels.
const (
	productionLog  = "production"
	developmentLog = "development"
)

// application info
const (
	serviceName = "players-api"
)

// loadSettings load all application parameters.
func loadSettings() *settings {
	webServerPort := loadIntEnvVar(webServerPortEnvVar)
	if webServerPort == 0 {
		webServerPort = 8080
	}

	newSettings := settings{
		repository:             loadRepositorySettings(),
		tracerServiceURL:       loadStringEnvVar(tracerServiceURL),
		logLevel:               loadStringEnvVar(logLevelEnvVar),
		passwordGenerationCost: loadIntEnvVar(passwordGenerationCostEnvVar),
		webServerPort:          webServerPort,
		grpcServerPort:         loadIntEnvVar(grpcServerPortEnvVar),
		timeoutToPublishSec:    loadIntEnvVar(timeoutToPublishSecEnvVar),
	}

	return &newSettings
}

// loadRepositorySettings load settings for player database.
func loadRepositorySettings() databaseSetup {
	return databaseSetup{
		dbName:   loadStringEnvVar(postgresDBNameEnvVar),
		host:     loadStringEnvVar(postgresHostEnvVar),
		player:   loadStringEnvVar(postgresPlayerEnvVar),
		password: loadStringEnvVar(postgresPasswordEnvVar),
		port:     loadIntEnvVar(postgresPortEnvVar),
	}
}

func loadIntEnvVar(key string) int {
	strValue := loadStringEnvVar(key)

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		slog.Error(
			"env var is not a number",
			slog.String("key", key),
			slog.String("value", strValue))

		return 0
	}

	return intValue
}

func loadStringEnvVar(key string) string {
	return os.Getenv(key)
}
