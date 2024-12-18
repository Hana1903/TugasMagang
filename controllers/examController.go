package controllers

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateExam(c *gin.Context) {
	var input struct {
		OrderID  int64   `json:"order_id"`
		PacketID int64   `json:"packet_id"`
		UserID   int64   `json:"user_id"`
		Score    float64 `json:"score"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the packet to get duration_exam
	var packet models.Packet
	if err := config.DB.First(&packet, input.PacketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found"})
		return
	}

	// Ambil durasi dari packet
	durationSeconds, err := strconv.Atoi(packet.DurationExam)
	if err != nil || durationSeconds <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format in packet", "details": err.Error()})
		return
	}

	// Set started_at to the current time and calculate ended_at
	startedAt := time.Now().In(time.Local)
	endedAt := startedAt.Add(time.Duration(durationSeconds) * time.Second)

	// Log untuk debugging
	log.Printf("StartedAt: %v, DurationSeconds: %v, EndedAt: %v", startedAt, durationSeconds, endedAt)

	// Create new Exam with calculated times
	exam := models.Exam{
		OrderID:   input.OrderID,
		PacketID:  input.PacketID,
		UserID:    input.UserID,
		Score:     input.Score,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		CreatedAt: time.Now(),
	}

	// Save the exam to the database
	if err := config.DB.Create(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam"})
		return
	}

	type Response struct {
		ID        int     `json:"id"`
		OrderID   int64   `json:"order_id"`
		PacketID  int64   `json:"packet_id"`
		UserID    int64   `json:"user_id"`
		Score     float64 `json:"score"`
		StartedAt string  `json:"started_at"`
		EndedAt   string  `json:"ended_at"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}
	
	// Membuat response struct
	response := Response{
		ID:        exam.ID,
		OrderID:   exam.OrderID,
		PacketID:  exam.PacketID,
		UserID:    exam.UserID,
		Score:     exam.Score,
		StartedAt: exam.StartedAt.Format("2006-01-02 15:04:05"),
		EndedAt:   exam.EndedAt.Format("2006-01-02 15:04:05"),
		CreatedAt: exam.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: exam.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// Kirim response
	c.JSON(http.StatusCreated, response)

}

func GetRemainingTime(c *gin.Context) {
	var exam models.Exam
	if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	remainingTime := time.Until(exam.EndedAt)
	if remainingTime < 0 {
		c.JSON(http.StatusOK, gin.H{"remaining_time": "0s"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"remaining_time": remainingTime.String()})
}

// Get all exams
func GetExams(c *gin.Context) {
	var exams []models.Exam
	if err := config.DB.Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exams", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exams)
}

// Get an exam by ID
func GetExamByID(c *gin.Context) {
	var exam models.Exam
	if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         exam.ID,
		"order_id":   exam.OrderID,
		"packet_id":  exam.PacketID,
		"user_id":    exam.UserID,
		"score":      exam.Score,
		"started_at": exam.StartedAt.Format("2006-01-02 15:04:05"),
		"ended_at":   exam.EndedAt.Format("2006-01-02 15:04:05"),
		"created_at": exam.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": exam.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

// Update an exam by ID
func UpdateExam(c *gin.Context) {
	var exam models.Exam
	if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	var input struct {
		OrderID   int64     `json:"order_id"`
		PacketID  int64     `json:"packet_id"`
		UserID    int64     `json:"user_id"`
		Score     float64   `json:"score"`
		StartedAt time.Time `json:"started_at"`
		EndedAt   time.Time `json:"ended_at"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Update hanya field yang tidak kosong
    if input.OrderID != 0 {
        exam.OrderID = input.OrderID
    }

    if input.PacketID != 0 {
        exam.PacketID = input.PacketID
    }

    if input.UserID != 0 {
        exam.UserID = input.UserID
    }

    if input.Score != 0 {
        exam.Score = input.Score
    }

    if !input.StartedAt.IsZero() {
        exam.StartedAt = input.StartedAt
    }

    if !input.EndedAt.IsZero() {
        exam.EndedAt = input.EndedAt
    }

    exam.UpdatedAt = time.Now()


	// Save changes to the database
	if err := config.DB.Save(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam", "details": err.Error()})
		return
	}

	type Response struct {
		ID        int     `json:"id"`
		OrderID   int64   `json:"order_id"`
		PacketID  int64   `json:"packet_id"`
		UserID    int64   `json:"user_id"`
		Score     float64 `json:"score"`
		StartedAt string  `json:"started_at"`
		EndedAt   string  `json:"ended_at"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	response := Response{
		ID:        exam.ID,
		OrderID:   exam.OrderID,
		PacketID:  exam.PacketID,
		UserID:    exam.UserID,
		Score:     exam.Score,
		StartedAt: exam.StartedAt.Format("2006-01-02 15:04:05"),
		EndedAt:   exam.EndedAt.Format("2006-01-02 15:04:05"),
		CreatedAt: exam.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: exam.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// Delete an exam by ID
func DeleteExam(c *gin.Context) {
	if err := config.DB.Delete(&models.Exam{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exam deleted successfully"})
}