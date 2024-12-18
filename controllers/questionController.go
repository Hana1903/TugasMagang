package controllers

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new question
func CreateQuestion(c *gin.Context) {
	var input struct {
		IDPackage     int64  `json:"id_package" binding:"required"`
		Question      string `json:"question" binding:"required"`
		Answer        string `json:"answer"`
		CorrectAnswer string `json:"correct_answer"`
		PacketID      int64  `json:"packet_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format answer to have newline between options
	formattedAnswer := strings.ReplaceAll(input.Answer, ", ", "\n")

	// Check if the answer is correct
	isCorrect := 0
	if input.Answer == input.CorrectAnswer {
		isCorrect = 1
	}

	question := models.Question{
		IDPackage:     input.IDPackage,
		Question:      input.Question,
		Answer:        formattedAnswer,
		CorrectAnswer: input.CorrectAnswer,
		PacketID:      input.PacketID,
		IsCorrect:     isCorrect,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

// 	var questions []models.Question
// 	config.DB.Find(&questions)
// 	c.JSON(http.StatusOK, questions)
// }

// Get all questions
func GetQuestions(c *gin.Context) {
	var questions []models.Question
	if err := config.DB.Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Periksa jika tidak ada data
	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found"})
		return
	}

	// Konversi setiap pertanyaan menggunakan ToResponse
	var responses []models.QuestionResponse
	for _, question := range questions {
		responses = append(responses, question.ToResponse())
	}

	// Gunakan ToResponse untuk memformat data
	c.JSON(http.StatusOK, responses)
}

// Get questions by package ID
func GetQuestionsByPackageID(c *gin.Context) {
	packageIDStr := c.Param("id_package")
	packageID, err := strconv.ParseInt(packageIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package ID"})
		return
	}

	var questions []models.Question
	if err := config.DB.Where("id_package = ?", packageID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No questions found for the given package ID"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// Get a question by ID
func GetQuestionByID(c *gin.Context) {
	var question models.Question
	if err := config.DB.Where("id = ?", c.Param("id")).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	c.JSON(http.StatusOK, question.ToResponse())
}

// Update a question by ID
func UpdateQuestion(c *gin.Context) {
	var question models.Question
	if err := config.DB.Where("id = ?", c.Param("id")).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var input struct {
		IDPackage     int64  `json:"id_package" binding:"required"`
		Question      string `json:"question" binding:"required"`
		Answer        string `json:"answer"`
		CorrectAnswer string `json:"correct_answer"`
		PacketID      int64  `json:"packet_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format answer to have newline between options
	formattedAnswer := strings.ReplaceAll(input.Answer, ", ", "\n")

	// Check if the answer is correct
	isCorrect := 0
	if input.Answer == input.CorrectAnswer {
		isCorrect = 1
	}

	// Update fields
	question.IDPackage = input.IDPackage
	question.Question = input.Question
	question.Answer = formattedAnswer
	question.CorrectAnswer = input.CorrectAnswer
	question.PacketID = input.PacketID
	question.IsCorrect = isCorrect
	question.UpdatedAt = time.Now()

	config.DB.Save(&question)
	c.JSON(http.StatusOK, question)
}

// Delete a question by ID
func DeleteQuestion(c *gin.Context) {
	var question models.Question
	if err := config.DB.Where("id = ?", c.Param("id")).First(&question).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	config.DB.Delete(&question)
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

// CalculateScore calculates the score based on user answers
func CalculateScore(c *gin.Context) {
	var input struct {
		Answers []struct {
			QuestionID int64  `json:"question_id" binding:"required"`
			Answer     string `json:"answer" binding:"required"`
		} `json:"answers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	correctAnswers := 1
	incorrectAnswers := 0
	totalQuestions := len(input.Answers)

	for _, userAnswer := range input.Answers {
		var question models.Question
		if err := config.DB.Where("id = ?", userAnswer.QuestionID).First(&question).Error; err == nil {
			if question.CorrectAnswer == userAnswer.Answer {
				correctAnswers++
				question.IsCorrect = 1
			} else {
				incorrectAnswers++
				question.IsCorrect = 0
			}
		}
	}

	scoreResult := models.ScoreResult{
		TotalQuestions:   totalQuestions,
		CorrectAnswers:   correctAnswers,
		IncorrectAnswers: incorrectAnswers,
	}

	c.JSON(http.StatusOK, scoreResult)
}
