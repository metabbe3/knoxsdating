// profile_repository.go
package repository

import (
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
)

type ProfileRepository interface {
	CreateProfile(profile *models.Profile) error
	GetProfileByID(profileID int) (*models.Profile, error)
	GetProfileByUserID(userID int) (*models.Profile, error)
	UpdateProfile(profile *models.Profile) error
	DeleteProfile(profile *models.Profile) error
}

type profileRepository struct {
	db    helpers.DatabaseHandler
	redis helpers.RedisHandler
}

func NewProfileRepository(db helpers.DatabaseHandler, redis helpers.RedisHandler) ProfileRepository {
	return &profileRepository{db: db, redis: redis}
}

func (r *profileRepository) CreateProfile(profile *models.Profile) error {
	result := r.db.Create(profile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *profileRepository) GetProfileByID(profileID int) (*models.Profile, error) {
	var profile models.Profile
	result := r.db.First(&profile, profileID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

func (r *profileRepository) GetProfileByUserID(userID int) (*models.Profile, error) {
	var profile models.Profile
	result := r.db.First(&profile).Where("UserID = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

func (r *profileRepository) UpdateProfile(profile *models.Profile) error {
	result := r.db.Save(profile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *profileRepository) DeleteProfile(profile *models.Profile) error {
	result := r.db.Delete(profile)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// NewUserRepositoryWithGormDBAndRedis creates a new ProfileRepository with GormDB and Redis
func NewProfileRepositoryWithGormDBAndRedis(db helpers.DatabaseHandler, redis helpers.RedisHandler) ProfileRepository {
	return NewProfileRepository(db, redis)
}
