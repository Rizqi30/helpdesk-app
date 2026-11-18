package services

import (
	"errors"
	"helpdesk-app/backend/models"
	"helpdesk-app/backend/repositories"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required,oneof=requester agent admin"`
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.UserRepo.FindAll()
}

func (s *UserService) UpdateUserRole(userID uint, input UpdateRoleInput) (*models.User, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		
		return nil, errors.New("Failed to get user")
	}

	if user.Role == input.Role {
		return user, nil
	}

	user.Role = input.Role

	if err := s.UserRepo.Update(user); err != nil {
		return nil, errors.New("Failed to update user role")
	}
	
	return user, nil
}