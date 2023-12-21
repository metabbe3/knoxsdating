// repository/notification_repository_test.go
package repository

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metabbe3/knoxsdating/pkg/helpers"
	"github.com/metabbe3/knoxsdating/pkg/helpers/mocks"
	"github.com/metabbe3/knoxsdating/pkg/models"
	"gorm.io/gorm"
	// Adjust the import path based on your project structure
	// Adjust the import path based on your project structure
)

func Test_notificationRepository_DeleteNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mockDB as a *mocks.MockDatabaseHandler
	var mockDB helpers.DatabaseHandler = &mocks.MockDatabaseHandler{}
	mockDBMock := mockDB.(*mocks.MockDatabaseHandler)

	mockRedis := mocks.MockRedisHandler{}                 // Create an instance of MockRedisHandler
	repo := NewNotificationRepository(mockDB, &mockRedis) // Pass both mockDB and mockRedis

	// Positive Test Case
	mockDBMock.DeleteFunc = func(value interface{}, conds ...interface{}) *gorm.DB {
		return &gorm.DB{} // You can customize the return value as needed
	}
	err := repo.DeleteNotification(&models.Notification{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Negative Test Case
	mockDBMock.DeleteFunc = func(value interface{}, conds ...interface{}) *gorm.DB {
		return &gorm.DB{Error: errors.New("mocked database error")}
	}
	err = repo.DeleteNotification(&models.Notification{})
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}
