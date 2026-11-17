package services

import (
	"errors"
	"helpdesk-app/backend/models"
	"helpdesk-app/backend/repositories"
	"helpdesk-app/backend/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
	JWTSecret string
	JWTExpiry time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
	}
}

type RegisterInput struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Register(input RegisterInput) (*models.User, string, error) {
	existingUser, err := s.UserRepo.FindByEmail(input.Email)
	if err == nil && existingUser.ID != 0 {
		return nil, "", errors.New("Email already registered")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", errors.New("Internal server error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New("Failed to hash password")
	}

	user := &models.User{
		Name: input.Name,
		Email: input.Email,
		PasswordHash: string(hashedPassword),
		Role: models.RoleRequester,
	}

	if err := s.UserRepo.Create(user); err != nil {
		return nil, "", errors.New("Failed to create user")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, s.JWTSecret, s.JWTExpiry)
	if err != nil {
		return nil, "", errors.New("Failed to generate token")
	}

	return user, token, nil
}

func (s *AuthService) Login(input LoginInput) (*models.User, string, error) {
	user, err := s.UserRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, "", errors.New("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", errors.New("Invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, s.JWTSecret, s.JWTExpiry)
	if err != nil {
		return nil, "", errors.New("Failed to generate token")
	}

	return user, token, nil
}

func (s *AuthService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	return user, nil
}