package handlers

import (
	"helpdesk-app/backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// GetAll - GET /api/admin/users
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

// UpdateRole - PATCH /api/admin/users/:id/role
func (h *UserHandler) UpdateRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID",
		})
		return
	}

	var input services.UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user, err := h.UserService.UpdateUserRole(uint(userID), input)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "User not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User role updated successfully",
		"data":    user,
	})
}