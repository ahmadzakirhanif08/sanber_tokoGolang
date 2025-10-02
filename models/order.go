package models

import (
	"time"
	"gorm.io/gorm"
)

type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	
	UserID       uint           `gorm:"not null" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID"` // Relasi ke User
	OrderDate    time.Time      `gorm:"default:now()" json:"order_date"`
	TotalAmount  float64        `gorm:"not null" json:"total_amount"`
	Status       string         `gorm:"not null;default:PENDING" json:"status"`
	
	Items        []OrderItem    // Relasi One-to-Many ke OrderItems
}