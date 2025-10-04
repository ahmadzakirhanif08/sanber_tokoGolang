package main

import (
	"fmt"
	"log"

	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/configs"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/models"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/ahmadzakirhanif08/sanber_tokoGolang.git/docs"
)


// @title Mini E-Commerce REST API - Golang Final Project
// @version 1.0
// @description Backend API untuk sistem e-commerce sederhana.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main(){
	
	//load .env
	err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file, assuming production environment.")
    }

	//connect to database
	configs.ConnectDatabase()

	//model migration
	configs.DB.AutoMigrate(
		&models.User{},
        &models.Product{},
        &models.Order{},
        &models.OrderItem{},
	)
	fmt.Println("Database migration success")

	//router setup
	router := gin.Default()

	//call router func
	routes.SetupRouter(router)

	//start server
	log.Println("Server starting on: ")
	router.SetTrustedProxies([]string{"127.0.0.1"})
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
