// package main
// import (
// 	"crud-ukom/config"
// 	"crud-ukom/models"
// 	"crud-ukom/routes"
// 	"log"
// 	"net/http"
// 	"github.com/rs/cors"
// )
// func main() {
// 	config.ConnectDB()
// 	config.DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Packet{}, &models.Exam{}, &models.Order{})
// 	router := routes.SetupRoutes()
// 	c := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"*"},                                                 // Mengizinkan semua origin, ganti dengan domain tertentu jika perlu
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},           // Metode yang diizinkan
// 		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"}, // Header yang diizinkan
// 		AllowCredentials: true,                                                          // Jika ingin mengizinkan cookies
// 	})
// 	handler := c.Handler(router)
// 	log.Fatal(http.ListenAndServe(":8080", handler))
// }

package main

import (
	"crud-ukom/config"
	"crud-ukom/models"
	"crud-ukom/routes"
)

func main() {
	config.ConnectDB()

	config.DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Packet{}, &models.Exam{}, &models.Order{})

	router := routes.SetupRoutes()
	router.Run("0.0.0.0:8080")
}
