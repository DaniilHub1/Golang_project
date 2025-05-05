package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID   	  uint
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Nickname  string
	PhotoPath string
	Role      string `gorm:"default:user"`
	Email     string // добавьте поле Email
}
