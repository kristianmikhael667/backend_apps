package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Contact struct {
	Provider    Provider  `json:"provider" gorm:"references:id"`
	ID          int16     `json:"id" gorm:"serial;primaryKey"`
	ContactId   uuid.UUID `json:"contact_id" gorm:"type:char(36);not_null;unique"`
	Phone       string    `json:"phone" gorm:"varchar;not null;unique"`
	ProviderId  int16     `json:"provider_id" gorm:"uuid;not null"`
	Status      int       `json:"status" gorm:"int;not null"`
	GanjilGenep string    `json:"ganjil_genep" gorm:"varchar;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (contact *Contact) BeforeSave(tx *gorm.DB) (err error) {
	if contact.ContactId == uuid.Nil {
		contact.ContactId = uuid.NewV4()
	}
	return nil
}
