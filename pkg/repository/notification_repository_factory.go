// repository/notification_repository_factory.go
package repository

import (
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"gorm.io/gorm"
)

func NewNotificationRepositoryWithConnection(dbHandler helpers.DatabaseHandler) NotificationRepository {
	return &notificationRepository{db: dbHandler}
}

func NewNotificationRepositoryWithGormDB(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: helpers.NewGormDBHandler(db)}
}
