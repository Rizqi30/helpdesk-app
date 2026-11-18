package services

import (
	"errors"
	"helpdesk-app/backend/models"
	"helpdesk-app/backend/repositories"

	"gorm.io/gorm"
)

type CommentService struct {
	CommentRepo *repositories.CommentRepository
	TicketRepo  *repositories.TicketRepository
}

func NewCommentService(commentRepo *repositories.CommentRepository, ticketRepo *repositories.TicketRepository) *CommentService {
	return &CommentService{
		CommentRepo: commentRepo,
		TicketRepo:  ticketRepo,
	}
}

type CreateCommentInput struct {
	Body string `json:"body" binding:"required,min=1"`
	IsInternal bool `json:"is_internal"`
}

func (s *CommentService) AddComment(ticketID uint, input CreateCommentInput, userID uint, userRole string) (*models.Comment, error) {
	ticket, err := s.TicketRepo.FindByID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Ticket not found")
		}
		return nil, errors.New("Failed to get ticket")
	}

	// Requester just comment at their own ticket
	if userRole == models.RoleRequester {
		if ticket.RequesterID != userID {
			return nil, errors.New("Ticket not found")
		}

		if input.IsInternal {
			return nil, errors.New("Requester cannot create internal notes")
		}
	}

	comment := &models.Comment{
		TicketID: ticketID,
		AuthorID: userID,
		Body: input.Body,
		IsInternal: input.IsInternal,
	}

	if err := s.CommentRepo.Create(comment); err != nil {
		return nil, errors.New("Failed to create comment")
	}

	// Reload with author for the full response
	comments, err := s.CommentRepo.FindByTicketID(ticketID, true)
	if err != nil {
		return nil, errors.New("Failed to load comment")
	}

	if len(comments) > 0 {
		return &comments[len(comments)-1], nil
	}

	return comment, nil
}

func (s *CommentService) GetComments(ticketID uint, userID uint, userRole string) ([]models.Comment, error) {
	ticket, err := s.TicketRepo.FindByID(ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Ticket not found")
		}
		return nil, errors.New("Failed to get ticket")
	}

	// Requester just see their own ticket
	if userRole == models.RoleRequester && ticket.RequesterID != userID {
		return nil, errors.New("Ticket not found")
	}

	// take the comment - requester didn't see internal notes
	includeInternal := userRole != models.RoleRequester

	return s.CommentRepo.FindByTicketID(ticketID, includeInternal)
}