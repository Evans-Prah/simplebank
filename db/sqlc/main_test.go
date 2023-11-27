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
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	testQueries = New(testDB)

	// Run tests
	exitCode := m.Run()

	// Close database connection
	if err := testDB.Close(); err != nil {
		log.Fatal("error closing the database connection", err)
	}

	os.Exit(exitCode)
}
