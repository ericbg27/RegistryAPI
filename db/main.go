package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CreateDBManager created the db manager using
func CreateDBManager(dbPath string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	db.AutoMigrate(&User{})

	return
}
