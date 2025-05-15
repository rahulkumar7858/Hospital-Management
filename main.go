package main

import (
	"hospital-management/config"
	"hospital-management/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	defer config.DB.Close()

	r := gin.Default()

	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)

	reception := r.Group("/reception")
	reception.Use(handlers.AuthMiddleware("receptionist"))
	{
		reception.POST("/patients", handlers.CreatePatient)
		reception.GET("/patients", handlers.GetPatients)
		reception.GET("/patients/:id", handlers.GetPatientByID)
		reception.PUT("/patients/:id", handlers.UpdatePatient)
		reception.DELETE("/patients/:id", handlers.DeletePatient)
	}

	doctor := r.Group("/doctor")
	doctor.Use(handlers.AuthMiddleware("doctor"))
	{
		doctor.GET("/patients", handlers.GetPatients)
		doctor.PUT("/patients/:id/diagnosis", handlers.UpdatePatientDiagnosis) // Fixed endpoint name
	}

	r.Static("/static", "./static")
	r.Run(":8081")
}
