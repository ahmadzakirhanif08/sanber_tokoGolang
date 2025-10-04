package handlers

import (
	"fmt"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/configs"
	"github.com/ahmadzakirhanif08/sanber_tokoGolang.git/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}

// @Summary Create a new order
// @Description Authenticated user creates an order, triggering database transaction and stock deduction.
// @Tags Orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body handlers.CreateOrderRequest true "Order Items List"
// @Success 201 {object} map[string]interface{} "Order created successfully"
// @Failure 400 {object} map[string]string "Bad Request (Invalid items, quantity, or insufficient stock)"
// @Failure 401 {object} map[string]string "Unauthorized (Missing or invalid token)"
// @Failure 500 {object} map[string]string "Internal Server Error (Transaction failed)"
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := configs.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}

	var totalAmount float64 = 0
	var orderItems []models.OrderItem

	order := models.Order{
		UserID: userID.(uint),
		Status: "PENDING",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order header"})
		return
	}

	for _, itemReq := range req.Items {
		var product models.Product

		if err := tx.First(&product, itemReq.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product ID %d not found", itemReq.ProductID)})
			return
		}
		if product.Stock < itemReq.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Insufficient stock for product %s (Available: %d)", product.Name, product.Stock)})
			return
		}

		subTotal := product.Price * float64(itemReq.Quantity)
		totalAmount += subTotal

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			SubTotal:  subTotal,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item"})
			return
		}

		orderItems = append(orderItems, orderItem)

		if err := tx.Model(&product).Update("stock", product.Stock-itemReq.Quantity).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deduct product stock"})
			return
		}
	}

	if err := tx.Model(&order).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update total amount"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order successfully created", "order_id": order.ID, "total": totalAmount})
}

func GetMyOrders(c *gin.Context) {
	var orders []models.Order
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	result := configs.DB.Preload("Items.Product").Where("user_id = ?", userID.(uint)).Find(&orders)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed fetching data"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "you have no order yet.", "data": []models.Order{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "data successfully fetched", "data": orders})
}
