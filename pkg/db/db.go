package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/memekas/ws-server/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) Init() error {

	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("PG_HOST"),     // host
		os.Getenv("PG_PORT"),     // port
		os.Getenv("PG_USER"),     // user
		os.Getenv("PG_DB"),       // dbname
		os.Getenv("PG_PASSWORD"), //pass
	)

	var err error
	db.con, err = gorm.Open("postgres", dbURI)
	if err != nil {
		return err
	}

	db.con.AutoMigrate(&Account{})
	return nil
}

func (db *DB) Get() *gorm.DB {
	return db.con
}

func (db *DB) Close() error {
	return db.con.Close()
}

func (db *DB) CreateUser(user *Account) error {
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
	tk := &auth.Token{}
	tkString, err := tk.Create(user.ID)
	if err != nil {
		return err
	}
	user.Token = tkString

	// Delete user password
	user.Password = ""

	return nil
}

func (db *DB) LoginUser(user *Account) error {
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
	tk := &auth.Token{}
	tkString, err := tk.Create(user.ID)
	if err != nil {
		return err
	}
	user.Token = tkString

	return nil
}
