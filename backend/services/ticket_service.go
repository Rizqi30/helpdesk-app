package services

import (
	"errors"
	"helpdesk-app/backend/models"
	"helpdesk-app/backend/repositories"

	"gorm.io/gorm"
)

type TicketService struct {
	TicketRepo *repositories.TicketRepository
	UserRepo *repositories.UserRepository
}

func NewTicketService(ticketRepo *repositories.TicketRepository, userRepo *repositories.UserRepository) *TicketService {
	return &TicketService{
		TicketRepo: ticketRepo,
		UserRepo: userRepo,
	}
}

type CreateTicketInput struct {
	Title string `json:"title" binding:"required,min=5,max=200"`
	Description string `json:"description" binding:"required,min=10"`
	Category string `json:"category" binding:"omitempty,max=50"`
	Priority string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
}

type UpdateTicketInput struct {
	Title *string `json:"title" binding:"omitempty,min=5,max=200"`
	Description *string `json:"description" binding:"omitempty,min=10"`
	Category *string `json:"category" binding:"omitempty,max=50"`
	Priority *string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status *string `json:"status" binding:"omitempty,oneof=open in_progress resolved closed"`
	AssigneeID *uint  `json:"assignee_id" binding:"omitempty"`
}

func (s *TicketService) CreateTicket(input CreateTicketInput, requesterID uint) (*models.Ticket, error) {
	priority := models.PriorityMedium
	if input.Priority != "" {
		priority = input.Priority
	}

	ticket := &models.Ticket{
		Title: input.Title,
		Description: input.Description,
		Category: input.Category,
		Priority: priority,
		Status: models.StatusOpen,
		RequesterID: requesterID,
	}

	if err := s.TicketRepo.Create(ticket); err != nil {
		return nil, errors.New("Failed to create ticket")
	}

	ticket, err := s.TicketRepo.FindByID(ticket.ID)
	if err != nil {
		return nil, errors.New("Failed to load ticket")
	}

	return ticket, nil
}

func (s *TicketService) GetTickets(filter repositories.TicketFilter, userID uint, userRole string) (*repositories.TicketListResult, error) {
	var requesterID uint
	if userRole == models.RoleRequester {
		requesterID = userID
	}

	return s.TicketRepo.FindAll(filter, requesterID)
}

func (s *TicketService) GetTicketByID(ticketID uint, userID uint, userRole string) (*models.Ticket, error) {
	ticket, err := s.TicketRepo.FindByID(ticketID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Ticket not found")
		}
		return nil, errors.New("Failed to get ticket")
	}

	if userRole == models.RoleRequester && ticket.RequesterID != userID {
		return nil, errors.New("Ticket not found")
	}

	return ticket, nil
}

func (s *TicketService) UpdateTicket(ticketID uint, input UpdateTicketInput, userID uint, userRole string) (*models.Ticket, error) {
	ticket, err := s.TicketRepo.FindByID(ticketID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Ticket not found")
		}
		return nil, errors.New("Failed to get ticket")
	}

	if userRole == models.RoleRequester {
		if ticket.RequesterID != userID {
			return nil, errors.New("Ticket not found")
		}

		if input.AssigneeID != nil {
			return nil, errors.New("Requester cannot assign tickets")
		}

		if input.Status != nil {
			return nil, errors.New("Requester cannot change ticket status")
		}
	}

	if input.AssigneeID != nil {
		assignee, err := s.UserRepo.FindByID(*input.AssigneeID)
		
		if err != nil {
			return nil, errors.New("Assignee not found")
		}

		if assignee.Role == models.RoleRequester {
			return nil, errors.New("Cannot assign ticket to a requester")
		}
	}

	if input.Title != nil {
		ticket.Title = *input.Title
	}
	if input.Description != nil {
		ticket.Description = *input.Description
	}
	if input.Category != nil {
		ticket.Category = *input.Category
	}
	if input.Priority != nil {
		ticket.Priority = *input.Priority
	}
	if input.Status != nil {
		ticket.Status = *input.Status
	}
	if input.AssigneeID != nil {
		ticket.AssigneeID = input.AssigneeID
	}

	if err := s.TicketRepo.Update(ticket); err != nil {
		return nil, errors.New("Failed to update ticket")
	}

	ticket, err = s.TicketRepo.FindByID(ticket.ID)
	if err != nil {
		return nil, errors.New("Failed to load ticket")
	}

	return ticket, nil
}