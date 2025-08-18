package main

import (
	"database/sql"
	"event-app/internal/database"
	"event-app/internal/env"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	// Muat .env bila ada
	_ = godotenv.Load()
	db, err := sql.Open("postgres", env.GetEnvString("PG_URI", ""))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
