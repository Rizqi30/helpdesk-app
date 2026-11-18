package repositories

import (
	"helpdesk-app/backend/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *CommentRepository) FindByTicketID(ticketID uint, includeInternal bool) ([]models.Comment, error) {
	var comments []models.Comment

	query := r.DB.Where("ticket_id = ?", ticketID)

	if !includeInternal {
		query = query.Where("is_internal = ?", false)
	}

	err := query.
		Order("created_at ASC").
		Preload("Author").
		Find(&comments).Error

	return comments, err
}