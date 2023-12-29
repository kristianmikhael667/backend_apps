package seeder

import (
	"backend_apps/models"
	"log"

	"time"

	"gorm.io/gorm"
)

func userSeeder(db *gorm.DB) {
	now := time.Now()

	var user = []models.User{
		{
			FullName:     "Super Admin",
			Email:        "super.admin@gmail.com",
			PasswordHash: "$2a$10$rfpS/jJ.a5J9seBM5sNPTeMQ0iVcAjoox3TDZqLE7omptkVQfaRwW", // 123abcABC!
			Status:       1,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Can't seeder data user, with error %v \n", err)
	}
	log.Println("Success seed data user")
}
