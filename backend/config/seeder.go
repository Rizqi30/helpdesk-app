package config

import (
	"helpdesk-app/backend/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB) {
	users := []models.User{
		{
			Name: "Admin",
			Email: "admin@helpdesk.com",
			PasswordHash: hashPassword("$$admin$$"),
			Role: models.RoleAdmin,
		},
		{
			Name: "Agent",
			Email: "agent@helpdesk.com",
			PasswordHash: hashPassword("$$agent$$"),
			Role: models.RoleAgent,
		},
		{
			Name: "User",
			Email: "user@helpdesk.com",
			PasswordHash: hashPassword("$$user$$"),
			Role: models.RoleRequester,
		},
	}

	for _, user := range users {
		result := db.Where("email = ?", user.Email).FirstOrCreate(&user)
		if result.Error != nil {
			log.Printf("Failed to seed user %s: %v", user.Email, result.Error)
			continue
		}

		if result.RowsAffected > 0 {
			log.Printf("Seeded user: %s (%s)", user.Email, user.Role)
		} else {
			log.Printf("User already exists: %s", user.Email)
		}
	}
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password: ", err)
	}
	return string(hashed)
}
