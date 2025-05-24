package models

import (
	"time"
)

type Message struct {
    ID         uint      `gorm:"primaryKey"`
    SenderID   uint      `gorm:"index"`
    ReceiverID uint      `gorm:"index"`
    Content    string
    ReplyToID  *uint     `gorm:"index;default:null"` // ID сообщения, на которое отвечаем
    CreatedAt  time.Time
}