package db_test

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (dbms *DBManagerSuite) TestCreateUser() {
	userMockRows := sqlmock.NewRows([]string{"id"}).AddRow("1")

	dbms.mock.ExpectBegin()
	dbms.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "users" ("full_name","phone","user_name","password") VALUES ($1,$2,$3,$4)`),
	).WithArgs(
		dbms.user.FullName,
		dbms.user.Phone,
		dbms.user.UserName,
		dbms.user.Password,
	).WillReturnRows(userMockRows)
	dbms.mock.ExpectCommit()

	user, err := dbms.manager.CreateUser(dbms.user)
	assert.NoError(dbms.T(), err)
	assert.Equal(dbms.T(), dbms.user, user)
}
