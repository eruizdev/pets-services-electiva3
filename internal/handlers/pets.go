package handlers

import (
	"net/http"
	"pets-service/internal/db"
	"pets-service/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllPets devuelve todas las mascotas (solo ADMIN y VET)
func GetAllPets(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, owner_id, name, species, breed, age FROM pets`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error consultando mascotas"})
		return
	}
	defer rows.Close()

	pets := []models.Pet{}
	for rows.Next() {
		var p models.Pet
		rows.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Species, &p.Breed, &p.Age)
		pets = append(pets, p)
	}
	c.JSON(http.StatusOK, pets)
}

// GetMyPets devuelve las mascotas del cliente autenticado
func GetMyPets(c *gin.Context) {
	// El gateway pasa el owner_id en el header
	ownerIDStr := c.GetHeader("X-User-Id")
	ownerID, _ := strconv.Atoi(ownerIDStr)

	rows, err := db.DB.Query(`SELECT id, owner_id, name, species, breed, age FROM pets WHERE owner_id = $1`, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error consultando mascotas"})
		return
	}
	defer rows.Close()

	pets := []models.Pet{}
	for rows.Next() {
		var p models.Pet
		rows.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Species, &p.Breed, &p.Age)
		pets = append(pets, p)
	}
	c.JSON(http.StatusOK, pets)
}

// CreatePet crea una nueva mascota asociada al cliente autenticado
func CreatePet(c *gin.Context) {
	ownerIDStr := c.GetHeader("X-User-Id")
	ownerID, _ := strconv.Atoi(ownerIDStr)

	var req models.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pet models.Pet
	query := `INSERT INTO pets (owner_id, name, species, breed, age) VALUES ($1,$2,$3,$4,$5) RETURNING id, owner_id, name, species, breed, age`
	err := db.DB.QueryRow(query, ownerID, req.Name, req.Species, req.Breed, req.Age).
		Scan(&pet.ID, &pet.OwnerID, &pet.Name, &pet.Species, &pet.Breed, &pet.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando mascota"})
		return
	}

	c.JSON(http.StatusCreated, pet)
}

// GetPetByID devuelve una mascota por su ID
func GetPetByID(c *gin.Context) {
	id := c.Param("id")

	var pet models.Pet
	err := db.DB.QueryRow(`SELECT id, owner_id, name, species, breed, age FROM pets WHERE id = $1`, id).
		Scan(&pet.ID, &pet.OwnerID, &pet.Name, &pet.Species, &pet.Breed, &pet.Age)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mascota no encontrada"})
		return
	}

	c.JSON(http.StatusOK, pet)
}

// UpdatePet actualiza una mascota
func UpdatePet(c *gin.Context) {
	id := c.Param("id")

	var req models.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.DB.Exec(
		`UPDATE pets SET name=$1, species=$2, breed=$3, age=$4 WHERE id=$5`,
		req.Name, req.Species, req.Breed, req.Age, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando mascota"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mascota actualizada"})
}

// DeletePet elimina una mascota
func DeletePet(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec(`DELETE FROM pets WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando mascota"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mascota eliminada"})
}
