package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	// import postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Init() *gorm.DB {
	var db *gorm.DB
	db, err := gorm.Open(os.Getenv("DB_TYPE"), fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	))

	if err != nil {
		log.Fatal(err, "Unable to connect database")
	}

	log.Info("Database Setup done...")

	AutoMigrate(db)
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.Debug().AutoMigrate(
		Contact{},
		StockWatch{},
		User{},
	)
}
