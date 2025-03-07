package baseModel

import "time"

type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `json:"status" gorm:"default:1"`
}

type BaseModelWithUUID struct {
	BaseModel
	UUID string `json:"uuid" gorm:"type:char(36);uniqueIndex"`
}
