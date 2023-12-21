package models

import "time"

type Message struct {
	MessageID      int       `gorm:"primaryKey"`
	SenderUserID   int       `gorm:"not null"`
	ReceiverUserID int       `gorm:"not null"`
	MessageContent string    `gorm:"type:text;not null"`
	Timestamp      time.Time `gorm:"type:timestamp"`
}

// Set the table name for the LocationHistory model
func (Message) TableName() string {
	return "Message"
}
