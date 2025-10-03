package routes

import (
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/handlers" 
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/middlewares"
	"github.com/gin-gonic/gin"
)


func SetupRouter(router *gin.Engine) {
	
	// Group API /api
	api := router.Group("/api")
	{
		// ROUTE Auth (Public)
		api.POST("/users/register", handlers.RegisterHandler)
		api.POST("/users/login", handlers.LoginHandler)

		// Product Group API
		products := api.Group("/products")
		{
			// CRUD Admin (Requires JWT + Admin Role)
			products.Use(middlewares.JWTAuthMiddleware(), middlewares.AdminAuthMiddleware()) 
			{
				products.POST("/", handlers.CreateProduct)
				products.PUT("/:id", handlers.UpdateProduct)
				products.DELETE("/:id", handlers.DeleteProduct)
			}
			
			// Read/View (Public/Guest Access)
			products.GET("/", handlers.GetAllProducts)
			products.GET("/:id", handlers.GetProductByID) 
		}

		// Order Group API (Requires JWT Login)
		orders := api.Group("/orders")
		orders.Use(middlewares.JWTAuthMiddleware())
		{
			orders.POST("/", handlers.CreateOrder)
			orders.GET("/", handlers.GetMyOrders)
		}
	}
}