package main

import (
	"helpdesk-app/backend/config"
	"helpdesk-app/backend/handlers"
	"helpdesk-app/backend/middlewares"
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

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)

	// Services
	jwtExpiry, err := time.ParseDuration(cfg.JWTExpiry)
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRY format", err)
	}

	authService := services.NewAuthService(userRepo, cfg.JWTSecret, jwtExpiry)
	ticketService := services.NewTicketService(ticketRepo, userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	ticketHandler := handlers.NewTicketHandler(ticketService)

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
	protected.Use(middlewares.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/me", authHandler.GetMe)

		protected.POST("/tickets", ticketHandler.Create)
		protected.GET("/tickets", ticketHandler.GetAll)
		protected.GET("/tickets/:id", ticketHandler.GetByID)
		protected.PATCH("/tickets/:id", ticketHandler.Update)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}