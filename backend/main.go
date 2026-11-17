package main

import (
	"helpdesk-app/backend/config"
	"helpdesk-app/backend/handlers"
	"helpdesk-app/backend/middleware"
	"helpdesk-app/backend/repositories"
	"helpdesk-app/backend/services"
	"log"
	"net/http"
	"time"

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

	userRepo := repositories.NewUserRepository(db)

	jwtExpiry, err := time.ParseDuration(cfg.JWTExpiry)
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRY format", err)
	}

	authServie := services.NewAuthService(userRepo, cfg.JWTSecret, jwtExpiry)
	authHandler := handlers.NewAuthHandler(authServie)

	// Router
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

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.GetMe)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}