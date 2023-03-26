package db

import (
	"gorm.io/gorm"
)

// User is the representation of an user in the database
type User struct {
	gorm.Model
	FullName   string
	Phone      string
	Password   string
	LoginToken string
}
