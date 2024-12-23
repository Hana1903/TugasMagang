package routes

import (
	"crud-ukom/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	// User routes
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	// Question routes
	r.GET("/questions", controllers.GetQuestions)
	r.GET("/questions/:id", controllers.GetQuestionByID)
	r.POST("/questions", controllers.CreateQuestion)
	r.PUT("/questions/:id", controllers.UpdateQuestion)
	r.DELETE("/questions/:id", controllers.DeleteQuestion)

	// Scoring route
	r.POST("/score", controllers.CalculateScore) // Calculate user score

	// Packet routes
	r.GET("/packets", controllers.GetPackets)
	r.GET("/packets-detail/:id", controllers.GetPacketByID)
	r.GET("/packets/:packet_id/questions", controllers.GetQuestionsByPacketID)
	r.GET("/packets-purchased/:id", controllers.GetPacketsByUser)
	r.POST("/packets", controllers.CreatePacket)
	r.PUT("/packets/:id", controllers.UpdatePacket)
	r.DELETE("/packets/:id", controllers.DeletePacket)

	//Hana
	// Exam Routes
	r.GET("/exams", controllers.GetExams)
	r.GET("/exams/:id", controllers.GetExamByID)
	r.GET("/exams/:id/remaining-time", controllers.GetRemainingTime)
	r.GET("/exams/packet", controllers.GetExamWithQuestions)
	r.POST("/exams", controllers.CreateExam)
	r.PUT("/exams/:id", controllers.UpdateExam)
	r.DELETE("/exams/:id", controllers.DeleteExam)

	// Order Routes
	r.GET("/orders", controllers.GetOrders)
	r.GET("/orders/:id", controllers.GetOrderByID)
	r.POST("/orders", controllers.CreateOrder)
	r.PUT("/orders/:id", controllers.UpdateOrder)
	r.DELETE("/orders/:id", controllers.DeleteOrder)


	return r
}
