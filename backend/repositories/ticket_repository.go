package repositories

import (
	"helpdesk-app/backend/models"
	"time"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

type TicketFilter struct {
	Status string
	Priority string
	AssigneeID *uint
	Search string
	SortBy string
	SortOrder string
	Page int
	Limit int
}

type TicketListResult struct {
	Tickets []models.Ticket `json:"tickets"`
	Total int64 `json:"total"`
	Page int `json:"page"`
	Limit int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

type AdminStats struct {
	TotalTickets int64 `json:"total_tickets"`
	OpenTickets int64 `json:"open_tickets"`
	InProgressTickets int64 `json:"in_progress_tickets"`
	ResolvedTickets int64 `json:"resolved_tickets"`
	ClosedTickets int64 `json:"closed_tickets"`
	OverdueTickets int64 `json:"overdue_tickets"`
	TotalUsers int64 `json:"total_users"`
}

func (r *TicketRepository) GetAdminStats() (*AdminStats, error) {
	var stats AdminStats

	if err := r.DB.Model(&models.Ticket{}).Count(&stats.TotalTickets).Error; err != nil {
		return nil, err
	}

	r.DB.Model(&models.Ticket{}).Where("status = ?", models.StatusOpen).Count(&stats.OpenTickets)
	r.DB.Model(&models.Ticket{}).Where("status = ?", models.StatusInProgress).Count(&stats.InProgressTickets)
	r.DB.Model(&models.Ticket{}).Where("status = ?", models.StatusResolved).Count(&stats.ResolvedTickets)
	r.DB.Model(&models.Ticket{}).Where("status = ?", models.StatusClosed).Count(&stats.ClosedTickets)

	twoHoursAgo := time.Now().Add(-models.SLAThreshold)
	r.DB.Model(&models.Ticket{}).
		Where("priority = ? AND status = ? AND created_at < ?", models.PriorityUrgent, models.StatusOpen, twoHoursAgo).
		Count(&stats.OverdueTickets)

	r.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	return &stats, nil
}

func (r *TicketRepository) FindAll(filter TicketFilter, requesterID uint) (*TicketListResult, error) {
	var tickets []models.Ticket
	var total int64

	query := r.DB.Model(&models.Ticket{})

	if requesterID > 0 {
		query = query.Where("requester_id = ?", requesterID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}

	if filter.Search != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	allowedSortBy := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"priority": true,
		"status": true,
		"title": true,
	}

	sortBy := "created_at"
	if allowedSortBy[filter.SortBy] {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filter.Page - 1) * filter.Limit

	err := query.
		Order(sortBy + " "+ sortOrder).
		Offset(offset).
		Limit(filter.Limit).
		Preload("Requester").
		Preload("Assignee").
		Find(&tickets).Error

	if err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.Limit
	if int(total)%filter.Limit > 0 {
		totalPages++
	}

	for i := range tickets {
		tickets[i].CalculateOverdue()
	}

	return &TicketListResult{
		Tickets: tickets,
		Total: total,
		Page: filter.Page,
		Limit: filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *TicketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.DB.
		Preload("Requester").
		Preload("Assignee").
		First(&ticket, id).Error

	if err == nil {
		ticket.CalculateOverdue()
	}
	
	return &ticket, err
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.DB.Create(ticket).Error
}

func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.DB.Save(ticket).Error
}