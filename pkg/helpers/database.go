// helpers/database_handler.go

package helpers

import (
	"time"

	"gorm.io/driver/postgres"
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
	Omit(columns ...string) *gorm.DB                  // Add Omit method
}

// GormDBHandler is the concrete implementation of DatabaseHandler for gorm.DB
type GormDBHandler struct {
	db *gorm.DB
}

// ConnectToDatabase connects to the PostgreSQL database with retry mechanism
func (g *GormDBHandler) ConnectToDatabase() (*gorm.DB, error) {
	dsn := "host=postgres user=knoxs password=knoxsdating dbname=knoxsdating port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	var db *gorm.DB
	var err error

	// Retry connecting to the database for a certain number of times
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		// Sleep for a short duration before retrying
		time.Sleep(2 * time.Second)
	}

	return nil, err
}

// NewDatabase creates a new GormDBHandler
func (g *GormDBHandler) NewDatabase() (*gorm.DB, error) {
	return g.ConnectToDatabase()
}

// Create implements the Create method from DatabaseHandler
func (g *GormDBHandler) Create(value interface{}) *gorm.DB {
	return g.db.Create(value)
}

// First implements the First method from DatabaseHandler
func (g *GormDBHandler) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.First(dest, conds...)
}

// Save implements the Save method from DatabaseHandler
func (g *GormDBHandler) Save(value interface{}) *gorm.DB {
	return g.db.Save(value)
}

// Delete implements the Delete method from DatabaseHandler
func (g *GormDBHandler) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Delete(value, conds...)
}

func (g *GormDBHandler) Where(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Where(query, args...)
}

func (g *GormDBHandler) Migrator() gorm.Migrator {
	return g.db.Migrator()
}

func (g *GormDBHandler) Table(name string) *gorm.DB {
	return g.db.Table(name)
}

func (g *GormDBHandler) Model(value interface{}) *gorm.DB {
	return g.db.Model(value)
}

// Find implements the Find method from DatabaseHandler
func (g *GormDBHandler) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Find(dest, conds...)
}

// Raw executes a raw SQL query
func (g *GormDBHandler) Raw(query string, values ...interface{}) *gorm.DB {
	return g.db.Raw(query, values...)
}

func (g *GormDBHandler) Omit(columns ...string) *gorm.DB {
	return g.db.Omit(columns...)
}

// NewGormDBHandler creates a new GormDBHandler
func NewGormDBHandler(db *gorm.DB) DatabaseHandler {
	return &GormDBHandler{db: db}
}
