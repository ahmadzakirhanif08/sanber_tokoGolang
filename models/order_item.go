package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID"` // Relasi ke Product
	Quantity  int     `gorm:"not null" json:"quantity"`
	SubTotal  float64 `gorm:"not null" json:"sub_total"`
}
