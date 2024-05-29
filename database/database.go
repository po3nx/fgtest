package database

import (
	"log"
	"os"

	"github.com/po3nx/fgtest/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)
func ConnectDb() {
    dsn := "sqlserver://admplanning:Planning2021@172.18.83.38/SQLSERVER?database=pmrsdev"
    db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
    }
	log.Println("connected")
	db.AutoMigrate(&models.Test{})
    DBConn = db
}
