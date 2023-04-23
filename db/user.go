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
	Admin      bool
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
		Admin:    false,
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

	result := dbManager.db.Omit("ID", "LoginToken").Limit(searchParams.Offset).Offset(searchOffset).Where("admin <> ?", true).Or("admin IS NULL").Find(&users)

	if err := result.Error; err != nil {
		return nil, err
	}

	return users, nil
}

type UpdateUserParams struct {
	ID         uint
	FullName   string
	Phone      string
	Password   string
	LoginToken string
}

func (dbManager *DBManager) UpdateUser(updateParams UpdateUserParams) error {
	result := dbManager.db.Model(&User{}).Where("id = ?", updateParams.ID).Select("full_name", "phone", "password", "login_token").Updates(User{
		FullName:   updateParams.FullName,
		Phone:      updateParams.Phone,
		Password:   updateParams.Password,
		LoginToken: updateParams.LoginToken,
	})

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &NotFoundError{
				object: "user",
			}
		}

		return err
	}

	return nil
}

func (dbManager *DBManager) DeleteUser(userName string) error {
	result := dbManager.db.Where("user_name = ?", userName).Delete(&User{})

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &NotFoundError{
				object: "user",
			}
		}

		return err
	}

	return nil
}
