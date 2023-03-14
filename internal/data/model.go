package data

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DatabaseItem interface {
	Validate(m *Manager) error
	SafeToDelete(m *Manager) error
}

type Item struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Link struct {
	Item
	Url         string `json:"link"`
	Description string `json:"description"`
	ReferenceID int    `json:"-"`
}

func (l *Link) Validate(*Manager) error {
	return errors.New("implement me")
}

func (l *Link) SafeToDelete(*Manager) error {
	return errors.New("implement me")
}

type Error struct {
	Message string `json:"message"`
}
