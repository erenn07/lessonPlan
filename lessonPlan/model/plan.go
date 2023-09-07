package model

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	State       string    `gorm:"not null;default:'processing'" json:"state" `
}
