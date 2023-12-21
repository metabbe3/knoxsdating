// models/swipe_history.go
package models

import (
	"time"
)

type SwipeHistory struct {
	SwipeHistoryEntityID int       `gorm:"column:SwipeHistoryEntityID;primaryKey;autoIncrement" json:"swipeHistoryEntityID"`
	SwiperUserID         int       `gorm:"column:SwiperUserID" json:"swiperUserID"`
	SwipedUserID         int       `gorm:"column:SwipedUserID" json:"swipedUserID"`
	SwipeDirection       string    `gorm:"column:SwipeDirection" json:"swipeDirection"`
	Timestamp            time.Time `gorm:"column:Timestamp" json:"timestamp"`
	RedoCount            int       `gorm:"column:RedoCount;default:1" json:"redoCount"`
	IsMatched            bool      `gorm:"column:IsMatched;default:false" json:"isMatched"`
}

// TableName specifies the table name for the SwipeHistory model
func (SwipeHistory) TableName() string {
	return "SwipeHistory"
}
