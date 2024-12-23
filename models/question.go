package models

import (
	"strings"
	"time"
)

type Question struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	PacketID      int64     `gorm:"type:bigint(20);default:null"`
	Question      string    `gorm:"type:text;not null"`
	Answer        string    `gorm:"type:varchar(255);default:null"`
	CorrectAnswer string    `gorm:"type:varchar(255);default:null"`
	IsCorrect     int       `gorm:"-"`
	CreatedAt     time.Time `gorm:"type:datetime(3);default:null"`
	UpdatedAt     time.Time `gorm:"type:datetime(3);default:null"`
}

// Custom output untuk JSON
type QuestionResponse struct {
	ID            int64    `json:"id"`
	PacketID      int64    `json:"packet_id"`
	Question      string   `json:"question"`
	Answer        []string `json:"answer"` // Array untuk format baris
	CorrectAnswer string   `json:"correct_answer"`
}

// Konversi Question ke QuestionResponse
func (q *Question) ToResponse() QuestionResponse {
	return QuestionResponse{
		ID:            q.ID,
		PacketID:      q.PacketID,
		Question:      q.Question,
		Answer:        strings.Split(q.Answer, "\n"), // Split jawaban berdasarkan newline
		CorrectAnswer: q.CorrectAnswer,
	}
}
