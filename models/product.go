package models

import (
	"time"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at`

	Name      string  `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Price     float64 `gorm:"not null" json:"price"`
	Stock     int     `gorm:"not null;default:0" json:"stock"`
}