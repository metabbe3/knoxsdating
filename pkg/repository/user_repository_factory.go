// repository/notification_repository_factory.go
package repository

import (
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"gorm.io/gorm"
)

func NewUserRepositoryWithConnection(dbHandler helpers.DatabaseHandler) UserRepository {
	return &userRepository{db: dbHandler}
}

func NewUserRepositoryWithGormDB(db *gorm.DB) UserRepository {
	return &userRepository{db: helpers.NewGormDBHandler(db)}
}
