package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	UserID    uint
	User      User
	PostID    uint
	Post      Post
	CreatedAt time.Time
}
