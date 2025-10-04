package handlers

import (
	"net/http"
	// "strconv"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/configs"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/models"
	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
}

// 1. CreateProduct: POST /api/products
func CreateProduct(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	}

	configs.DB.Create(&product)
	c.JSON(http.StatusCreated, gin.H{"message": "Produk created!", "data": product})
}

// 2. GetAllProducts: GET /api/products
func GetAllProducts(c *gin.Context) {
	var products []models.Product
	configs.DB.Find(&products)

	c.JSON(http.StatusOK, gin.H{"data": products})
}


// 3. GetProductByID: GET /api/products/:id
func GetProductByID(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := configs.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}


// 4. UpdateProduct: PUT /api/products/:id
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// 1. Check available product
	if err := configs.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 2. Bind new input
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Update data
	configs.DB.Model(&product).Updates(models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Produk updated!", "data": product})
}

// 5. DeleteProduct: DELETE /api/products/:id
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	// Cek produk ada
	if err := configs.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	configs.DB.Delete(&product)

	c.JSON(http.StatusOK, gin.H{"message": "Produk deleted!"})
}
