package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           int16     `json:"id" gorm:"serial;primaryKey"`
	Uid          uuid.UUID `json:"uid" gorm:"type:char(36);not null;unique"`
	FullName     string    `json:"full_name" gorm:"varchar"`
	Email        string    `json:"email" gorm:"varchar"`
	PasswordHash string    `json:"password_hash" gorm:"varchar"`
	Status       int8      `json:"status" gorm:"int2"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

func (user *User) BeforeSave(usr *gorm.DB) (err error) {
	if user.Uid == uuid.Nil {
		user.Uid = uuid.NewV4()
	}
	return nil
}

func (user *User) VerifyHash(otp, otpHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(otpHash), []byte(otp))
	if err == nil {
		return true
	}
	return false
}
