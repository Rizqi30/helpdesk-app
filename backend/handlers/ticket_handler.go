package handlers

import (
	"helpdesk-app/backend/repositories"
	"helpdesk-app/backend/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	TicketService *services.TicketService
}

func NewTicketHandler(ticketService *services.TicketService) *TicketHandler {
	return &TicketHandler{TicketService: ticketService}
}

// Create - POST /api/tickets
func (h *TicketHandler) Create(c *gin.Context) {
	var input services.CreateTicketInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	ticket, err := h.TicketService.CreateTicket(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"message": "Ticket created successfully",
		"data": ticket,
	})
}

// GetAll - GET /api/tickets
func (h *TicketHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	var assigneeID *uint
	if assigneeStr := c.Query("assignee_id"); assigneeStr != "" {
		id, err := strconv.ParseUint(assigneeStr, 10, 32)
		if err == nil {
			uid := uint(id)
			assigneeID = &uid
		}
	}

	filter := repositories.TicketFilter{
		Status: c.Query("status"),
		Priority: c.Query("priority"),
		AssigneeID: assigneeID,
		Search: c.Query("search"),
		SortBy: c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
		Page: page,
		Limit: limit,
	}

	result, err := h.TicketService.GetTickets(filter, userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": result,
	})
}

// GetByID - GET /api/tickets/:id
func (h *TicketHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(string)

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": "Invalid ticket ID",
		})
		return
	}

	ticket, err := h.TicketService.GetTicketByID(uint(ticketID), userID, userRole)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Ticket not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": ticket,
	})
}

// Update - PATCH /api/ticket/:id
func (h *TicketHandler) Update(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(string)

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	var input services.UpdateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	ticket, err := h.TicketService.UpdateTicket(uint(ticketID), input, userID, userRole)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "Ticket not found":
			status = http.StatusNotFound
		case "Requester cannot assign tickets",
			"Requester cannot change ticket priority",
			"Cannot assign ticket to a requester":
			status = http.StatusForbidden
		case "Assignee not found":
			status = http.StatusBadRequest
		default:
			if strings.Contains(err.Error(), "Cannot transition") ||
				strings.Contains(err.Error(), "is not allowed to transition") {
				status = http.StatusForbidden
			}
		}

		c.JSON(status, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Ticket updated successfully",
		"data": ticket,
	})
}