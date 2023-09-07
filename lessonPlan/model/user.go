package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int    `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	StudentNo int    `json:"studentno"`
}
