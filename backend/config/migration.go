package config

import (
	"helpdesk-app/backend/models"
	"log"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Ticket{},
		&models.Comment{},
	)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database migration completed successfully")
}