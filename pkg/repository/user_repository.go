// user_repository.go
package repository

import (
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(userID int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DoesUserWithEmailExist(email string) (bool, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
}

type userRepository struct {
	db    helpers.DatabaseHandler
	redis helpers.RedisHandler
}

func NewUserRepository(db helpers.DatabaseHandler, redis helpers.RedisHandler) UserRepository {
	return &userRepository{db: db, redis: redis}
}

func (r *userRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, `"User"."Email" = ?`, email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) DoesUserWithEmailExist(email string) (bool, error) {
	var count int64
	result := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) DeleteUser(user *models.User) error {
	result := r.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// NewUserRepositoryWithGormDBAndRedis creates a new UserRepository with GormDB and Redis
func NewUserRepositoryWithGormDBAndRedis(db helpers.DatabaseHandler, redis helpers.RedisHandler) UserRepository {
	return NewUserRepository(db, redis)
}
