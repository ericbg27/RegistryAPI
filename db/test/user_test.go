package db_test

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ericbg27/RegistryAPI/db"
	"github.com/stretchr/testify/assert"
)

func (dbms *DBManagerSuite) TestCreateUser() {
	userMockRows := sqlmock.NewRows([]string{"id"}).AddRow("0")

	dbms.mock.ExpectBegin()
	dbms.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","full_name","phone","user_name","password") VALUES ($1,$2,$3,$4,$5,$6,$7)`),
	).WithArgs(
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
		dbms.user.FullName,
		dbms.user.Phone,
		dbms.user.UserName,
		dbms.user.Password,
	).WillReturnRows(userMockRows)
	dbms.mock.ExpectCommit()

	userParams := db.CreateUserParams{
		FullName: dbms.user.FullName,
		Phone:    dbms.user.Phone,
		UserName: dbms.user.UserName,
		Password: dbms.user.Password,
	}

	user, err := dbms.manager.CreateUser(userParams)
	assert.NoError(dbms.T(), err)
	assert.Equal(dbms.T(), dbms.user.FullName, user.FullName)
	assert.Equal(dbms.T(), dbms.user.Phone, user.Phone)
	assert.Equal(dbms.T(), dbms.user.UserName, user.UserName)
	assert.Equal(dbms.T(), dbms.user.Password, user.Password)
}
