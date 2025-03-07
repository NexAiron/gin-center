package base

import (
	"time"
)

type Model interface {
	TableName() string
}
type BaseModel struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;index"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

func (m BaseModel) TableName() string {
	return ""
}

type PaginationRequest struct {
	Page     int `json:"page" binding:"required,min=1"`
	PageSize int `json:"page_size" binding:"required,min=1,max=100"`
}
type SortRequest struct {
	SortBy    string `json:"sort_by" binding:"omitempty"`
	SortOrder string `json:"sort_order" binding:"omitempty,oneof=asc desc"`
}
