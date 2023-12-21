// profile.go
package models

type Profile struct {
	ProfileID           int    `gorm:"column:ProfileID;primaryKey"`
	UserID              int    `gorm:"column:UserID;not null"`
	Photos              string `gorm:"column:Photos;type:jsonb"`
	AboutMe             string `gorm:"column:AboutMe;type:text"`
	Interests           string `gorm:"column:Interests;type:jsonb"`
	RelationshipGoals   string `gorm:"column:RelationshipGoals;type:text"`
	Height              int    `gorm:"column:Height;type:int"`
	Language            string `gorm:"column:Language;size:50"`
	ZodiacSign          string `gorm:"column:ZodiacSign;size:50"`
	EducationDetails    string `gorm:"column:EducationDetails;type:text"`
	SocialMediaAccounts string `gorm:"column:SocialMediaAccounts;type:jsonb"`
}

// Set the table name for the LocationHistory model
func (Profile) TableName() string {
	return "Profile"
}
