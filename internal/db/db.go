package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error abriendo DB: ", err)
	}

	// Reintentar conexion hasta 15 veces
	for i := 0; i < 15; i++ {
		if err = DB.Ping(); err == nil {
			break
		}
		log.Printf("DB no lista, reintento %d/15...", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatal("No se pudo conectar a DB: ", err)
	}

	createTables()
	log.Println("Pets DB conectada")
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS pets (
		id         SERIAL PRIMARY KEY,
		owner_id   INT          NOT NULL,
		name       VARCHAR(100) NOT NULL,
		species    VARCHAR(50)  NOT NULL,
		breed      VARCHAR(100),
		age        INT,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error creando tablas: ", err)
	}
}
