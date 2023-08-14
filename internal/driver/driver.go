package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection.
type DB struct {
	SQL *sql.DB
}

// dbConn is a global instance of DB.
var dbConn = &DB{}

const (
	maxOpenDbConn = 10
	maxIdleDbConn = 5
	maxDBLifetime = 5 * time.Minute
)

// NewDatabase creates a new database connection for the application.
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// ConnectSQL establishes a database connection pool for PostgreSQL.
func ConnectSQL(dsn string) (*DB, error) {
	// Create a new database connection.
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	// Set connection pool settings.
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	// Assign the database connection to the global dbConn instance.
	dbConn.SQL = d

	// Test the database connection.
	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB attempts to ping the database to check if the connection is valid.
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
