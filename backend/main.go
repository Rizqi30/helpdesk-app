package main

import (
	"helpdesk-app/backend/config"
	"helpdesk-app/backend/handlers"
	"helpdesk-app/backend/middlewares"
	"helpdesk-app/backend/models"
	"helpdesk-app/backend/repositories"
	"helpdesk-app/backend/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
	commentRepo := repositories.NewCommentRepository(db)

	// Services
	jwtExpiry, err := time.ParseDuration(cfg.JWTExpiry)
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRY format", err)
	}

	authService := services.NewAuthService(userRepo, cfg.JWTSecret, jwtExpiry)
	ticketService := services.NewTicketService(ticketRepo, userRepo, commentRepo)
	commentService := services.NewCommentService(commentRepo, ticketRepo)
	userService := services.NewUserService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userHandler := handlers.NewUserHandler(userService)

	// Router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"http://localhost",
			"http://localhost:80",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1",
		},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

		// Route tickets
		protected.POST("/tickets", ticketHandler.Create)
		protected.GET("/tickets", ticketHandler.GetAll)
		protected.GET("/tickets/:id", ticketHandler.GetByID)
		protected.PATCH("/tickets/:id", ticketHandler.Update)

		// Route comment
		protected.POST("/tickets/:id/comments", commentHandler.AddComment)
		protected.GET("/tickets/:id/comments", commentHandler.GetComments)
	}

	// Only admin routes
	admin:= r.Group("/api/users")
	admin.Use(middlewares.AuthMiddleware(cfg.JWTSecret), middlewares.RoleMiddleware(models.RoleAdmin))
	{
		admin.GET("", userHandler.GetAll)
		admin.GET("/stats", ticketHandler.GetStats)
		admin.PATCH("/:id/role", userHandler.UpdateRole)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}