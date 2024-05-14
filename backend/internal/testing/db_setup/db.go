package db_setup

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/panthershark/app/backend/db/connection"
)

func QuoteIdentifier(s string) string {
	return `"` + strings.Replace(s, `"`, `""`, -1) + `"`
}

type DBManager struct {
	muLock       sync.Mutex
	migrationDir string
	tmplName     string
	dbNames      []string
	connString   string
}

type DBOption func(dbm *DBManager)

func WithCustomTmplName(tmplName string) DBOption {
	return func(dbm *DBManager) {
		dbm.tmplName = tmplName
	}
}

func WithMigrationDir(migrationDir string) DBOption {
	return func(dbm *DBManager) {
		dbm.migrationDir = migrationDir
	}
}

func WithConnectionString(connString string) DBOption {
	return func(dbm *DBManager) {
		dbm.connString = connString
	}
}

// GetDatabaseManager: get a singleton instance of the test db manager.
func NewDatabaseManager(opts ...DBOption) *DBManager {
	dbm := DBManager{
		muLock:       sync.Mutex{},
		migrationDir: os.Getenv("DBMATE_MIGRATIONS_DIR"),
		tmplName:     "test_tmpl",
		dbNames:      []string{},
		connString:   os.Getenv("DATABASE_URL"),
	}

	for _, opt := range opts {
		opt(&dbm)
	}

	dbm.rmDatabases([]string{dbm.tmplName})
	dbm.createTemplateDatabase()

	return &dbm
}

func (dbm *DBManager) createTemplateDatabase() {
	conn, err := pgx.Connect(context.Background(), dbm.connString)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close(context.Background())

	var exists bool
	if err = conn.QueryRow(context.Background(), `SELECT EXISTS(SELECT * FROM pg_database WHERE datname = $1) as exists`, dbm.tmplName).Scan(&exists); err != nil {
		log.Panic(err)
	}

	if exists {
		log.Panicf("database %s already exists", dbm.tmplName)
	}

	dbUrl, err := url.Parse(dbm.connString)
	if err != nil {
		log.Panic(err)
	}

	// I suspect the connString needs to have dbo to create tables.
	// patching in the target dbname since dbmate says it can create it.
	dbUrl.Path = dbm.tmplName

	db := dbmate.New(dbUrl)
	db.MigrationsDir = []string{dbm.migrationDir}
	db.Strict = true

	if err = db.CreateAndMigrate(); err != nil {
		log.Panic(err)
	}

	if _, err = conn.Exec(context.Background(), fmt.Sprintf("ALTER DATABASE %s IS_TEMPLATE true", QuoteIdentifier(dbm.tmplName))); err != nil {
		log.Panic(err)
	}
}

func (dbm *DBManager) CreateTestDatabase() string {
	dbname := "test" + strings.ReplaceAll(uuid.New().String(), "-", "")
	conn, err := pgx.Connect(context.Background(), dbm.connString)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close(context.Background())

	if _, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", QuoteIdentifier(dbname), QuoteIdentifier(dbm.tmplName))); err != nil {
		log.Panic(err)
	}

	dbm.muLock.Lock()
	defer dbm.muLock.Unlock()

	dbm.dbNames = append(dbm.dbNames, dbname)

	dbUrl, err := url.Parse(dbm.connString)
	if err != nil {
		log.Panic(err)
	}

	dbUrl.Path = dbname

	return dbUrl.String()
}

func (dbm *DBManager) Cleanup() {
	dbm.rmDatabases(append(dbm.dbNames, dbm.tmplName))
}

func (dbm *DBManager) rmDatabases(dbNames []string) {
	connection.CloseDatabasePools()
	conn, err := pgx.Connect(context.Background(), dbm.connString)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close(context.Background())

	for _, dbName := range dbNames {
		exists := false
		conn.QueryRow(context.Background(), `SELECT EXISTS(SELECT * FROM pg_database WHERE datname = $1) as exists`, dbName).Scan(&exists)

		if exists {
			if _, err = conn.Exec(context.Background(), fmt.Sprintf("ALTER DATABASE %s IS_TEMPLATE false", QuoteIdentifier(dbName))); err != nil {
				log.Panic(err)
			}

			if _, err = conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS  %s", QuoteIdentifier(dbName))); err != nil {
				log.Panic(err)
			}
		}
	}

}
