package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Evans-Prah/simplebank/db/util"
	_ "github.com/lib/pq"
)



var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
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
