package model

import "time"

type Message struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	HandlerName string    `gorm:"type:varchar(255);not null" json:"handler_name"`
	Content     string    `gorm:"type:varchar(1024);not null" json:"content"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
