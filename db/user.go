package db

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName   string
	Phone      string `gorm:"unique"`
	UserName   string `gorm:"unique"`
	Password   string
	LoginToken string
}

type CreateUserParams struct {
	FullName string
	Phone    string
	UserName string
	Password string
}

func (dbManager *DBManager) CreateUser(userParams CreateUserParams) (*User, error) {
	user := &User{
		FullName: userParams.FullName,
		Phone:    userParams.Phone,
		UserName: userParams.UserName,
		Password: userParams.Password,
	}

	result := dbManager.db.Omit("LoginToken").Create(user)

	if err := result.Error; err != nil {
		if IsUniqueConstraintViolationError(err) {
			pqErr, _ := err.(*pq.Error)
			if pqErr.Column == "phone" {
				return nil, fmt.Errorf("An user with the provided phone number already exists")
			} else if pqErr.Column == "user_name" {
				return nil, fmt.Errorf("An user with the provided username already exists")
			}

			return nil, fmt.Errorf("An user with the provided information already exists")
		}

		return nil, err
	}

	return user, nil
}

func (dbManager *DBManager) GetUser(userName string) (*User, error) {
	var user User

	result := dbManager.db.Where("user_name = ?", userName).First(&user)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &NotFoundError{
				object: "user",
			}
		}

		return nil, err
	}

	return &user, nil
}

type GetUsersParams struct {
	PageIndex int
	Offset    int
}

func (dbManager *DBManager) GetUsers(searchParams GetUsersParams) ([]User, error) {
	var users []User

	searchOffset := searchParams.PageIndex * searchParams.Offset

	result := dbManager.db.Omit("ID", "LoginToken").Limit(searchParams.Offset).Offset(searchOffset).Find(&users)

	if err := result.Error; err != nil {
		return nil, err
	}

	return users, nil
}
