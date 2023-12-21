// repository/profile_repository_factory.go
package repository

import (
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"gorm.io/gorm"
)

func NewProfileRepositoryWithConnection(dbHandler helpers.DatabaseHandler) ProfileRepository {
	return &profileRepository{db: dbHandler}
}

func NewProfileRepositoryWithGormDB(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: helpers.NewGormDBHandler(db)}
}
