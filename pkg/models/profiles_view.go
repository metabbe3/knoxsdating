// models/profile_view.go
package models

import "time"

type ProfileView struct {
	ID           int       `gorm:"column:ProfileViewID;primaryKey"`
	ViewerUserID int       `gorm:"column:ViewerUserID;uniqueIndex:UniqueProfileView;not null"`
	ShownUserID  int       `gorm:"column:ShownUserID;uniqueIndex:UniqueProfileView;not null"`
	Timestamp    time.Time `gorm:"column:Timestamp;type:timestamp;uniqueIndex:UniqueViewPerDay"`
	DateOnly     time.Time `gorm:"column:DateOnly;type:date"`
}

// Set the table name for the ProfileView model
func (ProfileView) TableName() string {
	return "ProfileView"
}
