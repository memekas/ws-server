package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB - connection
type DB struct {
	con *gorm.DB
}

// Init - Open db connection and migrate Account struct
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

	db.con.Debug().AutoMigrate(&Account{})
	return nil
}

// Get - get DB connection
func (db *DB) Get() *gorm.DB {
	return db.con
}

// Close - close DB connection
func (db *DB) Close() error {
	return db.con.Close()
}
