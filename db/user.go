package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName   string
	Phone      string
	UserName   string
	Password   string
	LoginToken string
}

func (dbManager *DBManager) CreateUser(user *User) (*User, error) {
	result := dbManager.db.Omit("CreatedAt", "UpdatedAt", "DeletedAt", "LoginToken").Create(user)

	if err := result.Error; err != nil {
		return nil, err
	}

	return user, nil
}
