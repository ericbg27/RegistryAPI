package db

import "gorm.io/gorm"

type DBConnector interface {
	CreateUser(userParams CreateUserParams) (*User, error)
}

type DBManager struct {
	db *gorm.DB
}

// NewDBManager creates the db manager using the provided DB connection
func NewDBManager(db *gorm.DB) *DBManager {
	db.AutoMigrate(&User{})

	return &DBManager{
		db: db,
	}
}
