package repository

import (
	"github.com/metabbe3/knoxsdating/pkg/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db Database
}

type Database interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
}

func NewMessageRepository(db Database) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(message *models.Message) error {
	result := r.db.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MessageRepository) GetMessageByID(messageID int) (*models.Message, error) {
	var message models.Message
	result := r.db.First(&message, messageID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

func (r *MessageRepository) UpdateMessage(message *models.Message) error {
	result := r.db.Save(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MessageRepository) DeleteMessage(message *models.Message) error {
	result := r.db.Delete(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
