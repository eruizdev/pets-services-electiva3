package models

// Pet representa una mascota en la base de datos
type Pet struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Breed   string `json:"breed"`
	Age     int    `json:"age"`
}

// CreatePetRequest es lo que llega al crear una mascota
type CreatePetRequest struct {
	Name    string `json:"name"    binding:"required"`
	Species string `json:"species" binding:"required"`
	Breed   string `json:"breed"`
	Age     int    `json:"age"`
}
