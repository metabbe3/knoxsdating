// helpers/helpers.go
package helpers

import (
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

// DefaultDB is the default instance of DatabaseHandler
var DefaultDB DatabaseHandler = &GormDBHandler{}

// ConnectToDatabase connects to the database
func ConnectToDatabase() (*gorm.DB, error) {
	return DefaultDB.ConnectToDatabase()
}

// NewDatabase creates a new database connection
func NewDatabase() (*gorm.DB, error) {
	return DefaultDB.NewDatabase()
}

// ValidationError returns an error with the specified validation message.
func ValidationError(message string) error {
	return fmt.Errorf("validation error: %s", message)
}

// IsValidEmail checks if the given email is valid.
func IsValidEmail(email string) bool {
	// Regular expression for email validation
	// You can modify this regex pattern as per your requirements
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
