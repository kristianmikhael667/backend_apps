package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Provider struct {
	ID           int16     `json:"id" gorm:"serial;primaryKey"`
	ProviderId   uuid.UUID `json:"provider_id" gorm:"type:char(36);not_null;unique"`
	PrefixCode   string    `json:"prefix_code" gorm:"varchar;not null;unique"`
	NameProvider string    `json:"name_provider" gorm:"varchar;not null"`
	Status       int       `json:"status" gorm:"int;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Contact      []Contact `json:"contact"`
}

func (provider *Provider) BeforeSave(tx *gorm.DB) (err error) {
	if provider.ProviderId == uuid.Nil {
		provider.ProviderId = uuid.NewV4()
	}
	return nil
}
