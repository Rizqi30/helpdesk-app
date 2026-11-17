package main

import (
	"helpdesk-app/backend/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := config.ConnectDatabase(cfg)

	config.RunMigration(db)
	config.RunSeeder(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying DB: ", err)
	}
	defer sqlDB.Close()

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"message": "Database connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Server is running and database is connected",
		})
	})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}