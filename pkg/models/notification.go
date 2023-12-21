package models

import "time"

type Notification struct {
	NotificationID   int       `gorm:"column:NotificationID;primaryKey"`
	UserID           int       `gorm:"column:UserID;not null"`
	NotificationType string    `gorm:"column:NotificationType;size:50;not null"`
	Message          string    `gorm:"column:Message;type:text;not null"`
	Timestamp        time.Time `gorm:"column:Timestamp;type:timestamp"`
	IsRead           bool      `gorm:"column:IsRead;default:false"`
}

// Set the table name for the LocationHistory model
func (Notification) TableName() string {
	return "Notification"
}
