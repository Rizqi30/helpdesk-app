package handlers

import (
	"helpdesk-app/backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	CommentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}

// AddCOmment - POST /api/ticket/:id/comments
func (h *CommentHandler) AddComment(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(string)

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ticket ID",
		})
		return
	}

	var input services.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	comment, err := h.CommentService.AddComment(uint(ticketID), input, userID, userRole)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "Ticket not found":
			status = http.StatusNotFound
		case "Requester cannot create internal notes":
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Comment added successfully",
		"data":    comment,
	})
}

// GetComments - GET /api/tickets/:id/comments
func (h *CommentHandler) GetComments(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(string)

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ticket ID",
		})
		return
	}

	comments, err := h.CommentService.GetComments(uint(ticketID), userID, userRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Ticket not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   comments,
	})
}