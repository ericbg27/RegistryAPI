package db_test

import (
	"fmt"
	"regexp"
	"strconv"

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

func (dbms *DBManagerSuite) TestGetUser() {
	userMockRow := sqlmock.NewRows([]string{"id", "full_name", "phone", "user_name", "password"}).AddRow("0", dbms.user.FullName, dbms.user.Phone, dbms.user.UserName, dbms.user.Password)

	dbms.mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "users" WHERE user_name = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`),
	).WithArgs(
		dbms.user.UserName,
	).WillReturnRows(userMockRow)

	user, err := dbms.manager.GetUser(dbms.user.UserName)
	assert.NoError(dbms.T(), err)
	assert.Equal(dbms.T(), dbms.user.FullName, user.FullName)
	assert.Equal(dbms.T(), dbms.user.Phone, user.Phone)
	assert.Equal(dbms.T(), dbms.user.UserName, user.UserName)
	assert.Equal(dbms.T(), dbms.user.Password, user.Password)
}

func (dbms *DBManagerSuite) TestGetUsers() {
	userMockRows := sqlmock.NewRows([]string{"id", "full_name", "phone", "user_name", "password"})

	consideredNumUsers := numUsers - 3

	for i := 0; i < consideredNumUsers; i++ {
		userMockRows.AddRow(strconv.Itoa(i), dbms.users[i].FullName, dbms.users[i].Phone, dbms.users[i].UserName, dbms.users[i].Password)
	}

	dbms.mock.ExpectQuery(
		regexp.QuoteMeta(fmt.Sprintf(`SELECT "users"."created_at","users"."updated_at","users"."deleted_at","users"."full_name","users"."phone","users"."user_name","users"."password" FROM "users" WHERE "users"."deleted_at" IS NULL LIMIT %d OFFSET %d`, consideredNumUsers, 1*(consideredNumUsers))),
	).WithArgs().WillReturnRows(userMockRows)

	searchParams := db.GetUsersParams{
		PageIndex: 1,
		Offset:    numUsers - 3,
	}

	users, err := dbms.manager.GetUsers(searchParams)
	assert.NoError(dbms.T(), err)
	assert.Equal(dbms.T(), consideredNumUsers, len(users))
	for i := 0; i < consideredNumUsers; i++ {
		assert.Equal(dbms.T(), dbms.users[i].FullName, users[i].FullName)
		assert.Equal(dbms.T(), dbms.users[i].Phone, users[i].Phone)
		assert.Equal(dbms.T(), dbms.users[i].UserName, users[i].UserName)
		assert.Equal(dbms.T(), dbms.users[i].Password, users[i].Password)
	}
}
