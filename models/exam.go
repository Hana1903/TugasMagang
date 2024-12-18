package models

import (
	"time"
)

type Exam struct {
	ID        int       `gorm:"primaryKey"`
	OrderID   int64     `json:"order_id"`
	PacketID  int64     `json:"packet_id"`
	UserID    int64     `json:"user_id"`
	Score     float64   `json:"score"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
