package models

import "time"

const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"

	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusResolved   = "resolved"
	StatusClosed     = "closed"
)

type Ticket struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"type:varchar(200);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Category    string `json:"category" gorm:"type:varchar(50)"`
	Priority    string `json:"priority" gorm:"type:varchar(20);not null;default:'medium'"`
	Status      string `json:"status" gorm:"type:varchar(20);not null;default:'open'"`
	RequesterID uint   `json:"requester_id" gorm:"not null;index"`
	AssigneeID  *uint  `json:"assignee_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Requester User `json:"requester" gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	Assignee User `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID;constraint:OnDelete:SET NULL"`
}