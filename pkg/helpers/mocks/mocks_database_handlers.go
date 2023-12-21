// mocks/mock_database_handler.go
package mocks

import (
	"gorm.io/gorm"
)

// DatabaseHandler defines the methods for database operations
type DatabaseHandler interface {
	ConnectToDatabase() (*gorm.DB, error)
	NewDatabase() (*gorm.DB, error)
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Migrator() gorm.Migrator
	Table(name string) *gorm.DB
	Model(value interface{}) *gorm.DB // Add Model method
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Raw(query string, values ...interface{}) *gorm.DB // Add Raw method

}

// MockDatabaseHandler is a mock implementation of the DatabaseHandler interface
// MockDatabaseHandler is a mock implementation of the DatabaseHandler interface
type MockDatabaseHandler struct {
	ConnectToDatabaseFunc func() (*gorm.DB, error)
	NewDatabaseFunc       func() (*gorm.DB, error)
	CreateFunc            func(value interface{}) *gorm.DB
	FirstFunc             func(dest interface{}, conds ...interface{}) *gorm.DB
	SaveFunc              func(value interface{}) *gorm.DB
	DeleteFunc            func(value interface{}, conds ...interface{}) *gorm.DB
	WhereFunc             func(query interface{}, args ...interface{}) *gorm.DB
	MigratorFunc          func() gorm.Migrator
	TableFunc             func(name string) *gorm.DB
	ModelFunc             func(interface{}) *gorm.DB
	FindFunc              func(dest interface{}, conds ...interface{}) *gorm.DB
	RawFunc               func(query string, values ...interface{}) *gorm.DB // Add Raw method
	OmitFunc              func(columns ...string) *gorm.DB                   // Add Omit method
}

// ConnectToDatabase implements the ConnectToDatabase method from DatabaseHandler
func (m *MockDatabaseHandler) ConnectToDatabase() (*gorm.DB, error) {
	if m.ConnectToDatabaseFunc != nil {
		return m.ConnectToDatabaseFunc()
	}
	return nil, nil
}

// NewDatabase implements the NewDatabase method from DatabaseHandler
func (m *MockDatabaseHandler) NewDatabase() (*gorm.DB, error) {
	if m.NewDatabaseFunc != nil {
		return m.NewDatabaseFunc()
	}
	return nil, nil
}

// Create implements the Create method from DatabaseHandler
func (m *MockDatabaseHandler) Create(value interface{}) *gorm.DB {
	if m.CreateFunc != nil {
		return m.CreateFunc(value)
	}
	return nil
}

// First implements the First method from DatabaseHandler
func (m *MockDatabaseHandler) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.FirstFunc != nil {
		return m.FirstFunc(dest, conds...)
	}
	return nil
}

// Save implements the Save method from DatabaseHandler
func (m *MockDatabaseHandler) Save(value interface{}) *gorm.DB {
	if m.SaveFunc != nil {
		return m.SaveFunc(value)
	}
	return nil
}

// Delete implements the Delete method from DatabaseHandler
func (m *MockDatabaseHandler) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(value, conds...)
	}
	return nil
}

// Where implements the Where method from DatabaseHandler
func (m *MockDatabaseHandler) Where(query interface{}, args ...interface{}) *gorm.DB {
	if m.WhereFunc != nil {
		return m.WhereFunc(query, args...)
	}
	return nil
}

// Migrator implements the Migrator method from DatabaseHandler
func (m *MockDatabaseHandler) Migrator() gorm.Migrator {
	if m.MigratorFunc != nil {
		return m.MigratorFunc()
	}
	return nil
}

// Table implements the Table method from DatabaseHandler
func (m *MockDatabaseHandler) Table(name string) *gorm.DB {
	if m.TableFunc != nil {
		return m.TableFunc(name)
	}
	return nil
}

// Model implements the Model method from DatabaseHandler
func (m *MockDatabaseHandler) Model(value interface{}) *gorm.DB {
	if m.ModelFunc != nil {
		return m.ModelFunc(value)
	}
	return nil
}

// Find implements the Find method from DatabaseHandler
func (m *MockDatabaseHandler) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.FindFunc != nil {
		return m.FindFunc(dest, conds...)
	}
	return nil
}

// Omit implements the Omit method from DatabaseHandler
func (m *MockDatabaseHandler) Omit(columns ...string) *gorm.DB {
	if m.OmitFunc != nil {
		return m.OmitFunc(columns...)
	}
	return nil
}

// Raw implements the Raw method from DatabaseHandler
func (m *MockDatabaseHandler) Raw(query string, values ...interface{}) *gorm.DB {
	if m.RawFunc != nil {
		return m.RawFunc(query, values...)
	}
	return nil
}
