package main

import (
	"fmt"
	"log"

	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/configs"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/models"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/handlers"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" 
)

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

	//group API
	api := router.Group("/api")
	{
        //ROUTE Auth
        api.POST("/users/register", handlers.RegisterHandler)
        api.POST("/users/login", handlers.LoginHandler)

		//product group api
		products := api.Group("/products")
        {
            products.Use(middlewares.BasicAuthMiddleware()) 
            {
                products.POST("/", handlers.CreateProduct)
                products.PUT("/:id", handlers.UpdateProduct)
                products.DELETE("/:id", handlers.DeleteProduct)
            }
            products.GET("/", handlers.GetAllProducts)
            products.GET("/:id", handlers.GetProductByID) 
        }

		//order group api
		orders := api.Group("/orders")
        orders.Use(middlewares.JWTAuthMiddleware())
        {
            orders.POST("/", handlers.CreateOrder)
            orders.GET("/", handlers.GetMyOrders)
        }
        
    }

	//start server
	log.Println("Server starting on: ")
	router.SetTrustedProxies([]string{"127.0.0.1"})
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
