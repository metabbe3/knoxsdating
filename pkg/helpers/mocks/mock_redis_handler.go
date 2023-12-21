// helpers/mocks/mock_redis_handler.go
package mocks

import (
	"time"
)

// MockRedisHandler is a mock implementation of the RedisHandler interface
type MockRedisHandler struct {
	// Implement methods of RedisHandler as needed for testing
	SetFunc    func(key string, value interface{}, expiration time.Duration) error
	GetFunc    func(key string, dest interface{}) error
	DeleteFunc func(key string) error
}

// Set implements the Set method from RedisHandler
func (m *MockRedisHandler) Set(key string, value interface{}, expiration time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(key, value, expiration)
	}
	return nil
}

// Get implements the Get method from RedisHandler
func (m *MockRedisHandler) Get(key string, dest interface{}) error {
	if m.GetFunc != nil {
		return m.GetFunc(key, dest)
	}
	return nil
}

// Delete implements the Delete method from RedisHandler
func (m *MockRedisHandler) Delete(key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(key)
	}
	return nil
}
