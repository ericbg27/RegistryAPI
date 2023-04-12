package db_test

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ericbg27/RegistryAPI/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const numUsers = 5

type DBManagerSuite struct {
	suite.Suite
	conn *sql.DB
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	manager *db.DBManager
	user    *db.User
	users   []*db.User
}

func (dbms *DBManagerSuite) SetupSuite() {
	var err error

	dbms.conn, dbms.mock, err = sqlmock.New()
	assert.NoError(dbms.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 dbms.conn,
		PreferSimpleProtocol: true,
	})

	dbms.DB, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(dbms.T(), err)

	dbms.manager = db.NewDBManager(dbms.DB)
	assert.IsType(dbms.T(), &db.DBManager{}, dbms.manager)

	dbms.user = &db.User{
		FullName: "Test User",
		Phone:    "99999999",
		UserName: "test",
		Password: "secret",
	}

	for i := 0; i < numUsers; i++ {
		userFullName := "Test User " + strconv.Itoa(i)
		userPhone := "9999999" + strconv.Itoa(i)
		userUserName := "test" + strconv.Itoa(i)
		userPassword := "secret" + strconv.Itoa(i)

		dbms.users = append(dbms.users, &db.User{
			FullName: userFullName,
			Phone:    userPhone,
			UserName: userUserName,
			Password: userPassword,
		})
	}
}

func (dbms *DBManagerSuite) AfterTest(_, _ string) {
	assert.NoError(dbms.T(), dbms.mock.ExpectationsWereMet())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DBManagerSuite))
}
