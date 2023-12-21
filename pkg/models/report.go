package models

import "time"

type Report struct {
	ReportID       int       `gorm:"primaryKey"`
	ReporterUserID int       `gorm:"not null"`
	ReportedUserID int       `gorm:"not null"`
	ReportContent  string    `gorm:"type:text;not null"`
	Timestamp      time.Time `gorm:"type:timestamp"`
}

// Set the table name for the LocationHistory model
func (Report) TableName() string {
	return "Report"
}
