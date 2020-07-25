package models

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Account - user struct
type Account struct {
	gorm.Model
	Email    string `json:"email" gorm:"UNIQUE"`
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}

// Create - create new user in *DB
func (user *Account) Create(db *DB) error {

	// Get hash from user.Password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPass)

	// Add user to db
	db.Get().Create(user)
	if user.ID <= 0 {
		return fmt.Errorf("Failed to create new user")
	}

	// Create jwt
	tk := &Token{}
	tkString, err := tk.Create(user.ID)
	if err != nil {
		return err
	}
	user.Token = tkString

	// Delete user password
	user.Password = ""

	return nil
}

// Login - login user
func (user *Account) Login(db *DB) error {
	password := user.Password

	err := db.Get().Table("accounts").Where("email = ?", user.Email).First(user).Error
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}

	user.Password = ""

	// Create jwt
	tk := &Token{}
	tkString, err := tk.Create(user.ID)
	if err != nil {
		return err
	}
	user.Token = tkString

	return nil
}

// Token JWT
type Token struct {
	jwt.StandardClaims
	UserID uint
}

// Create new JWT token
func (tk *Token) Create(UserID uint) (string, error) {
	tk.UserID = UserID
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Decrypt jwt token from string
func (tk *Token) Decrypt(tkString string) error {
	_, err := jwt.ParseWithClaims(tkString, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Notification that sends to users
type Notification struct {
	ToUser   uint   `json:"toUser"`
	FromUser uint   `json:-`
	Msg      string `json:"msg"`
}
