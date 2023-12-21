package models

import "time"

type User struct {
	UserID             int       `gorm:"column:UserID;primaryKey"`
	Username           string    `gorm:"column:Username;not null;unique"`
	Email              string    `gorm:"column:Email;not null;unique"`
	Password           string    `gorm:"column:Password;not null"`
	VerificationStatus bool      `gorm:"column:VerificationStatus;default:false"`
	VerificationBadge  bool      `gorm:"column:VerificationBadge;default:false"`
	PremiumStatus      string    `gorm:"column:PremiumStatus;default:'Free';not null"`
	PremiumStartDate   time.Time `gorm:"column:PremiumStartDate;type:timestamp"`
	PremiumEndDate     time.Time `gorm:"column:PremiumEndDate;type:timestamp"`
	Gender             string    `gorm:"column:Gender;size:20"`
	Company            string    `gorm:"column:Company;size:255"`
	School             string    `gorm:"column:School;size:255"`
	JobTitle           string    `gorm:"column:JobTitle;size:255"`
	VerifiedBadge      bool      `gorm:"column:VerifiedBadge;default:false"`
}

// Set the table name for the LocationHistory model
func (User) TableName() string {
	return "User"
}
