package config

import (
	"fmt"
	"lessonPlan/model"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize connection object.
func DatabaseInit() {
	godotenv.Load()
	dbhost := os.Getenv("MYSQL_HOST")
	dbuser := os.Getenv("MYSQL_USER")
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DBNAME")
	dbport := os.Getenv("MYSQL_PORT")

	c := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbuser, dbpassword, dbhost, dbport)
	database, err := gorm.Open(mysql.Open(c), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + ";")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpassword, dbhost, dbport, dbname)
	db, err := gorm.Open(mysql.Open(connection), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	AutoMigrate(db)
	fmt.Println("Succesfull")
}

func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(&model.Plan{})
	connection.Debug().AutoMigrate(&model.User{})
}
