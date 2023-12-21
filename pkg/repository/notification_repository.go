// notification_repository.go
package repository

import (
	"fmt"

	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/models"
)

type NotificationRepository interface {
	CreateNotification(notification *models.Notification) error
	GetNotificationByID(notificationID int) (*models.Notification, error)
	UpdateNotification(notification *models.Notification) error
	DeleteNotification(notification *models.Notification) error

	// New methods for Redis
	SaveNotificationToRedis(notification *models.Notification) error
	GetNotificationFromRedis(notificationID int) (*models.Notification, error)
}

type notificationRepository struct {
	db    helpers.DatabaseHandler
	redis helpers.RedisHandler
}

func NewNotificationRepository(db helpers.DatabaseHandler, redis helpers.RedisHandler) NotificationRepository {
	return &notificationRepository{db: db, redis: redis}
}

func (r *notificationRepository) CreateNotification(notification *models.Notification) error {
	result := r.db.Create(notification)
	if result.Error != nil {
		return result.Error
	}

	// Save to Redis after successful database creation
	if err := r.SaveNotificationToRedis(notification); err != nil {
		// Handle Redis error (log, return an error, etc.)
		fmt.Printf("Error saving to Redis: %v\n", err)
	}

	return nil
}

func (r *notificationRepository) GetNotificationByID(notificationID int) (*models.Notification, error) {
	// Try to get from Redis first
	notification, err := r.GetNotificationFromRedis(notificationID)
	if err == nil {
		return notification, nil
	}

	// If not found in Redis, fetch from the database
	var dbNotification models.Notification
	result := r.db.First(&dbNotification, notificationID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Save to Redis for future requests
	if err := r.SaveNotificationToRedis(&dbNotification); err != nil {
		// Handle Redis error (log, return an error, etc.)
		fmt.Printf("Error saving to Redis: %v\n", err)
	}

	return &dbNotification, nil
}

func (r *notificationRepository) UpdateNotification(notification *models.Notification) error {
	result := r.db.Save(notification)
	if result.Error != nil {
		return result.Error
	}

	// Update in Redis after successful database update
	if err := r.SaveNotificationToRedis(notification); err != nil {
		// Handle Redis error (log, return an error, etc.)
		fmt.Printf("Error saving to Redis: %v\n", err)
	}

	return nil
}

func (r *notificationRepository) DeleteNotification(notification *models.Notification) error {
	result := r.db.Delete(notification)
	if result.Error != nil {
		return result.Error
	}

	// Delete from Redis after successful database delete
	if err := r.redis.Delete(fmt.Sprintf("notification:%d", notification.NotificationID)); err != nil {
		// Handle Redis error (log, return an error, etc.)
		fmt.Printf("Error deleting from Redis: %v\n", err)
	}

	return nil
}

func (r *notificationRepository) SaveNotificationToRedis(notification *models.Notification) error {
	// Save to Redis with a key (you can use notificationID as the key)
	return r.redis.Set(fmt.Sprintf("notification:%d", notification.NotificationID), notification, 0)
}

func (r *notificationRepository) GetNotificationFromRedis(notificationID int) (*models.Notification, error) {
	// Try to get data from Redis
	var notification models.Notification
	err := r.redis.Get(fmt.Sprintf("notification:%d", notificationID), &notification)
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

// NewNotificationRepositoryWithGormDBAndRedis creates a new NotificationRepository with GormDB and Redis
func NewNotificationRepositoryWithGormDBAndRedis(db helpers.DatabaseHandler, redis helpers.RedisHandler) NotificationRepository {
	return NewNotificationRepository(db, redis)
}
