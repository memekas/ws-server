package db

import "github.com/jinzhu/gorm"

type DB struct {
	con *gorm.DB
}

type Account struct {
	gorm.Model
	Email    string `json:"email" gorm:"UNIQUE"`
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}
