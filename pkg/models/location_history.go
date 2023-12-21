// models/location_history.go

package models

import "time"

// LocationHistory represents the location history of a user.
type LocationHistory struct {
	LocationID int       `json:"locationID" db:"LocationID" gorm:"column:LocationID;primaryKey"`
	UserID     int       `json:"userID" db:"UserID" gorm:"column:UserID"`
	Latitude   float64   `json:"latitude" db:"Latitude" gorm:"column:Latitude"`
	Longitude  float64   `json:"longitude" db:"Longitude" gorm:"column:Longitude"`
	Timestamp  time.Time `json:"timestamp" db:"Timestamp" gorm:"column:Timestamp"`
}

// Set the table name for the LocationHistory model
func (LocationHistory) TableName() string {
	return "Locationhistory"
}
