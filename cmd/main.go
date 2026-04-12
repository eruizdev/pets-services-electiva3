package main

import (
	"log"
	"pets-service/internal/db"
	"pets-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	router := gin.Default()

	pets := router.Group("/pets")
	{
		// El gateway ya validó el JWT, aquí solo procesamos
		pets.GET("/", handlers.GetAllPets)      // ADMIN, VET
		pets.GET("/my", handlers.GetMyPets)     // CLIENT
		pets.POST("/", handlers.CreatePet)      // CLIENT
		pets.GET("/:id", handlers.GetPetByID)   // todos
		pets.PUT("/:id", handlers.UpdatePet)    // CLIENT, ADMIN
		pets.DELETE("/:id", handlers.DeletePet) // ADMIN
	}

	log.Println("🐾 Pets Service corriendo en puerto 8082")
	router.Run(":8082")
}
