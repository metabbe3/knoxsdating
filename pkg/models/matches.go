// models/match.go
package models

import (
	"time"
)

type Match struct {
	MatchID   int       `gorm:"column:MatchID" json:"matchID"`
	UserID1   int       `gorm:"column:UserID1" json:"userID1"`
	UserID2   int       `gorm:"column:UserID2" json:"userID2"`
	Timestamp time.Time `gorm:"column:Timestamp" json:"timestamp"`
}

// TableName specifies the table name for the SwipeHistory model
func (Match) TableName() string {
	return "Match"
}
