package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	testQueries = New(conn)

	// Run tests
	exitCode := m.Run()

	// Close database connection
	if err := conn.Close(); err != nil {
		log.Fatal("error closing the database connection", err)
	}

	os.Exit(exitCode)
}
