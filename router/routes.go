package routes

import (
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/handlers" 
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/middlewares"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func SetupRouter(router *gin.Engine) {

	//swagger router
	//url http://localhost:8080/swagger/index.html

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
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