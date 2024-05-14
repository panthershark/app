package db_setup_test

import (
	"context"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/panthershark/app/backend/internal/testing/db_setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type dbSetupSuite struct {
	suite.Suite
	connString string
}

func (suite *dbSetupSuite) SetupSuite() {
	suite.connString = os.Getenv("DATABASE_URL")
}

func TestDatabaseSetup(t *testing.T) {
	err := godotenv.Load("../../../test.env")
	if err != nil {
		log.Panic(err)
	}

	suite.Run(t, new(dbSetupSuite))
}

// ----- Start tests

func (suite *dbSetupSuite) TestDatabaseManager() {
	t := suite.T()

	tmplDB := "db_setup_test"
	var testDB string

	dbm := db_setup.NewDatabaseManager(
		db_setup.WithConnectionString(suite.connString),
		db_setup.WithCustomTmplName(tmplDB),
	)

	testConnString := dbm.CreateTestDatabase()
	dbUrl, err := url.Parse(testConnString)
	assert.NoError(t, err)
	testDB = dbUrl.Path[1:]

	conn, err := pgx.Connect(context.Background(), suite.connString)
	assert.NoError(t, err)
	defer conn.Close(context.Background())

	tmplDBExists := false
	conn.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT * FROM pg_database WHERE datname = $1) as exists",
		tmplDB,
	).Scan(&tmplDBExists)
	assert.True(t, tmplDBExists, "template db should exist")

	testDBExists := false
	conn.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT * FROM pg_database WHERE datname = $1) as exists",
		testDB,
	).Scan(&testDBExists)
	assert.True(t, testDBExists, "test db should exist")

	// Test Cleanup
	dbm.Cleanup()

	tmplDBWasDeleted := false
	conn.QueryRow(
		context.Background(),
		"SELECT NOT EXISTS(SELECT * FROM pg_database WHERE datname = $1) as exists",
		tmplDB,
	).Scan(&tmplDBWasDeleted)
	assert.True(t, tmplDBWasDeleted, "template db should be cleaned up")

	testDBWasDeleted := false
	conn.QueryRow(
		context.Background(),
		"SELECT NOT EXISTS(SELECT * FROM pg_database WHERE datname = $1) as exists",
		testDB,
	).Scan(&testDBWasDeleted)
	assert.True(t, testDBWasDeleted, "template db should be cleaned up")

}
