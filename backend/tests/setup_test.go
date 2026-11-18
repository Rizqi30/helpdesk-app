package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"helpdesk-app/backend/handlers"
	"helpdesk-app/backend/middlewares"
	"helpdesk-app/backend/models"
	"helpdesk-app/backend/repositories"
	"helpdesk-app/backend/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testRouter *gin.Engine
)

func TestMain(m *testing.M) {
	setupTestDB()
	setupRouter()

	code := m.Run()

	cleanupTestDB()

	os.Exit(code)
}

func setupTestDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getTestEnv("DB_HOST", "localhost"),
		getTestEnv("DB_PORT", "5432"),
		getTestEnv("DB_USER", "postgres"),
		getTestEnv("DB_PASSWORD", "admin"),
		getTestEnv("DB_NAME", "helpdesk_test"),
	)

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database: ", err)
	}

	// Migrate tabel
	testDB.AutoMigrate(&models.User{}, &models.Ticket{}, &models.Comment{})
}

func setupRouter() {
	gin.SetMode(gin.TestMode)

	userRepo := repositories.NewUserRepository(testDB)
	ticketRepo := repositories.NewTicketRepository(testDB)
	commentRepo := repositories.NewCommentRepository(testDB)

	jwtSecret := getTestEnv("JWT_SECRET", "test-secret-key")
	jwtExpiry := 24 * time.Hour

	authService := services.NewAuthService(userRepo, jwtSecret, jwtExpiry)
	ticketService := services.NewTicketService(ticketRepo, userRepo, commentRepo)
	commentService := services.NewCommentService(commentRepo, ticketRepo)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	ticketHandler := handlers.NewTicketHandler(ticketService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		protected.GET("/me", authHandler.GetMe)
		protected.POST("/tickets", ticketHandler.Create)
		protected.GET("/tickets", ticketHandler.GetAll)
		protected.GET("/tickets/:id", ticketHandler.GetByID)
		protected.PATCH("/tickets/:id", ticketHandler.Update)
		protected.POST("/tickets/:id/comments", commentHandler.AddComment)
		protected.GET("/tickets/:id/comments", commentHandler.GetComments)
	}

	admin := r.Group("/api/users")
	admin.Use(middlewares.AuthMiddleware(jwtSecret), middlewares.RoleMiddleware(models.RoleAdmin))
	{
		admin.GET("", userHandler.GetAll)
		admin.PATCH("/:id/role", userHandler.UpdateRole)
	}

	testRouter = r
}

// cleanupTestDB membersihkan semua data dari tabel test.
func cleanupTestDB() {
	sqlDB, _ := testDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func cleanTables() {
	testDB.Exec("DELETE FROM comments")
	testDB.Exec("DELETE FROM tickets")
	testDB.Exec("DELETE FROM users")
}

// Helper Functions 
func getTestEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func createTestUser(name, email, password, role string) *models.User {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         role,
	}
	testDB.Create(user)
	return user
}

func loginAndGetToken(email, password string) string {
	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return ""
	}
	token, _ := data["token"].(string)
	return token
}

func performRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}