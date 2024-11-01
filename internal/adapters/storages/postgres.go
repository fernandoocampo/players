package storages

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // just to load drivers
)

// Parameters contains postgres information.
type Parameters struct {
	Host     string
	Player   string
	Password string
	DBName   string
	Port     int
}

// NewPostgresClient creates a new postgresql connection.
func NewPostgresClient(parameters Parameters) (*sql.DB, error) {
	psqlInfo := buildPostgresqlConnection(parameters)

	pgsqlconn, err := connectToDatabase(psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pgsqlconn, nil
}

func buildPostgresqlConnection(dbParameters Parameters) string {
	return fmt.Sprintf(
		"host=%s port=%d player=%s password=%s dbname=%s sslmode=disable",
		dbParameters.Host,
		dbParameters.Port,
		dbParameters.Player,
		dbParameters.Password,
		dbParameters.DBName,
	)
}

// connectToDatabase creates a connection to postgresql database based on given client parameters.
func connectToDatabase(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// ensure connection calling ping method
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseConnection(conn *sql.DB) error {
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("unable to close connection: %w", err)
	}

	return nil
}
