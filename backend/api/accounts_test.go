package api_test

import (
	"context"
	"log"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/panthershark/app/backend/api"
	"github.com/panthershark/app/backend/db/connection"
	"github.com/panthershark/app/backend/db/dbc"
	"github.com/panthershark/app/backend/internal/testing/db_setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type dataSetupFn func(c *pgxpool.Conn) error

type accountsSuite struct {
	suite.Suite
	dbm       *db_setup.DBManager
	setupTest func(setup dataSetupFn) connection.Pool
}

func (suite *accountsSuite) SetupSuite() {
	suite.dbm = db_setup.NewDatabaseManager(db_setup.WithCustomTmplName("accounts_test"))

	suite.setupTest = func(setupTestData dataSetupFn) connection.Pool {
		connString := suite.dbm.CreateTestDatabase()
		pool := connection.NewPgxPool(connString)

		if err := pool.AcquireFunc(context.Background(), setupTestData); err != nil {
			log.Panicf("TEST SETUP FAILED. err=%v", err)
		}

		return pool
	}
}
func (suite *accountsSuite) TearDownSuite() {
	suite.dbm.Cleanup()
}

func TestAccountsIntegration(t *testing.T) {
	err := godotenv.Load("../test.env")
	if err != nil {
		log.Panic(err)
	}

	suite.Run(t, new(accountsSuite))
}

func (suite *accountsSuite) TestGetUserByEmail() {
	t := suite.T()

	t.Run("ServiceAccount", func(t *testing.T) {
		pool := suite.setupTest(func(c *pgxpool.Conn) error {
			return c.QueryRow(context.Background(), `INSERT INTO public.accounts (id,email) VALUES($1,$2) RETURNING id`, "95e3c393-0299-40f6-8896-d375d602e52f", "robot@svc.com").Scan(&uuid.UUID{})
		})

		var got *api.Account
		var err error
		pool.AcquireFunc(context.Background(), func(c *pgxpool.Conn) error {
			testApi := api.NewAccountsApi(dbc.New(c), nil)
			got, err = testApi.GetUserByEmail("robot@svc.com")
			return nil
		})

		assert.NoError(t, err)
		cupaloy.SnapshotT(t, got)
	})

	t.Run("UserAccount", func(t *testing.T) {
		pool := suite.setupTest(func(c *pgxpool.Conn) error {
			if err := c.QueryRow(context.Background(), `INSERT INTO public.accounts (id,email) VALUES($1,$2) RETURNING id`, "a8f474b0-400e-492c-8541-0026d3b9a047", "hulk@hogan.com").Scan(&uuid.UUID{}); err != nil {
				return err
			}

			return c.QueryRow(context.Background(), `INSERT INTO public.person (id,account_id,first_name,last_name) VALUES($1,$2,$3,$4) RETURNING id`, "04e7a054-31d4-48e8-b313-5c38031636b2", "a8f474b0-400e-492c-8541-0026d3b9a047", "Hulk", "Hogan").Scan(&uuid.UUID{})
		})

		var got *api.Account
		var err error
		pool.AcquireFunc(context.Background(), func(c *pgxpool.Conn) error {
			testApi := api.NewAccountsApi(dbc.New(c), nil)
			got, err = testApi.GetUserByEmail("hulk@hogan.com")
			return nil
		})

		assert.NoError(t, err)
		cupaloy.SnapshotT(t, got)
	})

	t.Run("NotFound", func(t *testing.T) {
		pool := suite.setupTest(func(c *pgxpool.Conn) error {
			return nil
		})

		var got *api.Account
		var err error
		pool.AcquireFunc(context.Background(), func(c *pgxpool.Conn) error {
			testApi := api.NewAccountsApi(dbc.New(c), nil)
			got, err = testApi.GetUserByEmail("nope@nope.com")
			return nil
		})

		assert.Error(t, err, "not found")
		assert.Nil(t, got)
	})
}
