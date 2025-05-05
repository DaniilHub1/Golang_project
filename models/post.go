package models

import (
	"time"
)

type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"not null"`
	UserID    uint      
	User      User      // связь с юзером
	CreatedAt time.Time
    Title string 
	IsAdmin   bool  
}
