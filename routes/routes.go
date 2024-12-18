package routes

import (
	"crud-ukom/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// User routes
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	// Question routes
	r.POST("/questions", controllers.CreateQuestion)
	r.GET("/questions", controllers.GetQuestions)
	r.GET("/questions/:id", controllers.GetQuestionByID)                          
	r.PUT("/questions/:id", controllers.UpdateQuestion)                           
	r.DELETE("/questions/:id", controllers.DeleteQuestion)                        
	r.GET("/packages/:id_package/questions", controllers.GetQuestionsByPackageID) 

	// Scoring route
	r.POST("/score", controllers.CalculateScore) // Calculate user score

	// Packet routes
	r.GET("/packets", controllers.GetPackets)
	r.GET("/packets/:id", controllers.GetPacketByID)
	r.POST("/packets", controllers.CreatePacket)
	r.PUT("/packets/:id", controllers.UpdatePacket)
	r.DELETE("/packets/:id", controllers.DeletePacket)


	//Hana
	// Exam Routes
	r.POST("/exams", controllers.CreateExam)
	r.GET("/exams", controllers.GetExams)
	r.GET("/exams/:id", controllers.GetExamByID)
	r.PUT("/exams/:id", controllers.UpdateExam)
	r.DELETE("/exams/:id", controllers.DeleteExam)
	r.GET("/exams/:id/remaining-time", controllers.GetRemainingTime)
    
	// Order Routes
	r.POST("/orders", controllers.CreateOrder)
	r.GET("/orders", controllers.GetOrders)
	r.GET("/orders/:id", controllers.GetOrderByID)
	r.PUT("/orders/:id", controllers.UpdateOrder)
	r.DELETE("/orders/:id", controllers.DeleteOrder)

	return r
}
