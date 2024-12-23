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

// Get exams by packet ID
func GetExamsByPacket(c *gin.Context) {
	// Ambil parameter package_id dari query string
	PacketID := c.Query("packet_id")
	if PacketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "packet_id is required"})
		return
	}

	// Konversi package_id menjadi integer
	id, err := strconv.Atoi(PacketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid packet_id"})
		return
	}

	// Query untuk mengambil semua exam berdasarkan package_id
	var exams []models.Exam
	if err := config.DB.Where("packet_id = ?", id).Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve exams", "details": err.Error()})
		return
	}

	// Cek apakah data ditemukan
	if len(exams) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no exams found for the given packet_id"})
		return
	}

	// Berikan response dalam format JSON
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

func GetExamWithQuestions(c *gin.Context) {
	// Ambil ID exam dari parameter
	examID := c.Param("id")

	// Cari exam berdasarkan ID
	var exam models.Exam
	if err := config.DB.First(&exam, examID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	// Cari semua pertanyaan berdasarkan packet_id dari exam
	var questions []models.Question
	if err := config.DB.Where("packet_id = ?", exam.PacketID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	// Format data output
	var questionResponses []models.QuestionResponse
	for _, question := range questions {
		questionResponses = append(questionResponses, question.ToResponse())
	}

	// Response final
    c.JSON(http.StatusOK, gin.H{
        "exam": exam,
        "questions": questionResponses,
    })
}

// Update an exam by ID
func UpdateExam(c *gin.Context) {
    var exam models.Exam
    if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ujian tidak ditemukan"})
        return
    }

    // Cek waktu ujian
    if exam.EndedAt.Before(time.Now()) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Ujian sudah selesai, tidak bisa diedit"})
        return
    }

    var input struct {
        OrderID int64 `json:"order_id"`
        PacketID int64 `json:"packet_id"`
		UserID int64 `json:"user_id"`
        Score float64 `json:"score"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
        return
    }

    // Perbarui hanya field yang diperlukan
    exam.OrderID = input.OrderID
    exam.PacketID = input.PacketID
	exam.UserID = input.UserID
    exam.Score = input.Score
    exam.UpdatedAt = time.Now()

    if err := config.DB.Save(&exam).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui ujian"})
        return
    }

    // Kirim respons
    type Response struct {
        ID int `json:"id"`
        OrderID int64 `json:"order_id"`
        PacketID int64 `json:"packet_id"`
		UserID int64 `json:"user_id"`
        Score float64 `json:"score"`
        StartedAt string `json:"started_at"`
        EndedAt string `json:"ended_at"`
        CreatedAt string `json:"created_at"`
        UpdatedAt string `json:"updated_at"`
    }

    response := Response{
        ID: exam.ID,
        OrderID: exam.OrderID,
        PacketID: exam.PacketID,
		UserID: exam.UserID,
        Score: exam.Score,
        StartedAt: exam.StartedAt.Format("2006-01-02 15:04:05"),
        EndedAt: exam.EndedAt.Format("2006-01-02 15:04:05"),
        CreatedAt: exam.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: exam.UpdatedAt.Format("2006-01-02 15:04:05"),
    }

    c.JSON(http.StatusOK, gin.H{
		"id":         response.ID,
		"order_id":   response.OrderID,
		"packet_id":  response.PacketID,
		"user_id":    response.UserID,
		"score":      response.Score,
		"started_at": response.StartedAt,
		"ended_at":   response.EndedAt,
		"created_at": response.CreatedAt,
		"updated_at": response.UpdatedAt,
        "message": "Ujian berhasil diperbarui",
    })
}


// Delete an exam by ID
func DeleteExam(c *gin.Context) {
    var exam models.Exam
    if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ujian tidak ditemukan"})
        return
    }

    // Cek waktu ujian
    if exam.EndedAt.Before(time.Now()) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Ujian sudah selesai, tidak bisa dihapus"})
        return
    }

    if err := config.DB.Delete(&exam).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus ujian"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Ujian berhasil dihapus",
    })
}

