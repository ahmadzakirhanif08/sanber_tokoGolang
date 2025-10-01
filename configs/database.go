package configs

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	var database string

	if os.Getenv("DATABASE_URL") != "" {
		database = os.Getenv("DATABASE_URL")
	} else {
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbport := os.Getenv("DB_PORT")

		database = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", 
		dbHost, dbUser, dbPassword, dbName, dbport)
	}

	db, err := gorm.Open(postgres.Open(database), &gorm.Config{})
	if err != nil {
		panic("Database Connection Failed!")
	}
	DB = db
	if os.Getenv("DATABSE_URL") == ""{
		fmt.Println("Connection to local database estabilished!")
	} else {
		fmt.Println("Connection to prod database estabilished!")
	}
	

}
