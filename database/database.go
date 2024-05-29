package database

import (
	"log"
	"os"
	"fmt"
	"strconv"
	"github.com/po3nx/fgtest/models"
	"github.com/po3nx/fgtest/config"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)
func ConnectDb() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(err)
	}
    dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
	config.Config("DB_USER"),
	config.Config("DB_PASSWORD"),
	config.Config("DB_HOST"),
	port,
	config.Config("DB_NAME"))
    db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
    }
	log.Println("connected")
	db.AutoMigrate(&models.Test{})
    DBConn = db
}
