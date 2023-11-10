package models

import (
	"errors"
	token "food-siam-si/food-siam-si-user/utils"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	Customer UserType = "Customer"
	Owner    UserType = "Owner"
)

type User struct {
	UserId   uint     `gorm:"primary_key;auto_increment" json:"userId"`
	Username string   `gorm:"size:255;not null;unique" json:"username"`
	Password string   `gorm:"size:255;not null;" json:"password"`
	Email    string   `gorm:"size:100;not null;unique" json:"email"`
	UserType UserType `sql:"type:ENUM('Customer','Owner');" json:"userType"`
}

func (u *User) SaveUser() (*User, error) {

	u.BeforeSave()

	err := DB.Create(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil

}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return "randomtoken", nil

	token, err := token.GenerateToken(u.UserId)

	if err != nil {
		return "", err
	}

	return token, nil

}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}
	u.Password = ""

	return u, nil

}
