package models

import "time"

type Comment struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	TicketID   uint   `json:"ticket_id" gorm:"not null;index"`
	AuthorID   uint   `json:"author_id" gorm:"not null;index"`
	Body       string `json:"body" gorm:"type:text;not null"`
	IsInternal bool   `json:"is_internal" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`

	Ticket Ticket `json:"-" gorm:"foreignKey:TicketID;constraint:OnDelete:CASCADE"`
	Author User `json:"author" gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
}