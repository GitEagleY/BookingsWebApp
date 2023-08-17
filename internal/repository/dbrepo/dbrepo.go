package dbrepo

import (
	"database/sql"

	"github.com/GitEagleY/BookingsWebApp/internal/config"
	repository "github.com/GitEagleY/BookingsWebApp/internal/repository"
)

// postgresDBRepo is a struct that represents a PostgreSQL database repository.
type postgresDBRepo struct {
	App *config.AppConfig // App contains the application's configuration settings.
	DB  *sql.DB           // DB is the PostgreSQL database connection.
}
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates and returns a new instance of postgresDBRepo.
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,    // Initialize the App field with the provided AppConfig.
		DB:  conn, // Initialize the DB field with the provided PostgreSQL database connection.
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a, // Initialize the App field with the provided AppConfig.
		//DB:  conn, // Initialize the DB field with the provided PostgreSQL database connection.
	}
}
