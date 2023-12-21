package models

import "time"

type SystemLog struct {
	LogID      int       `gorm:"primaryKey"`
	LogType    string    `gorm:"size:20;not null"`
	LogMessage string    `gorm:"type:text;not null"`
	Timestamp  time.Time `gorm:"type:timestamp"`
}

// Set the table name for the LocationHistory model
func (SystemLog) TableName() string {
	return "SystemLog"
}
