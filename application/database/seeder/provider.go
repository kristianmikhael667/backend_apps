package seeder

import (
	"backend_apps/models"
	"log"

	"time"

	"gorm.io/gorm"
)

func providerSeeder(db *gorm.DB) {
	now := time.Now()

	var provider = []models.Provider{
		{
			PrefixCode:   "081",
			NameProvider: "Telkomsel",
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			PrefixCode:   "085",
			NameProvider: "Indosat",
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			PrefixCode:   "087",
			NameProvider: "XL",
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			PrefixCode:   "089",
			NameProvider: "Tri",
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			PrefixCode:   "088",
			NameProvider: "Smartfreen",
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			PrefixCode:   "083",
			NameProvider: "Axis",
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
	if err := db.Create(&provider).Error; err != nil {
		log.Printf("Can't seeder data provider, with error %v \n", err)
	}
	log.Println("Success seed data provider")
}
