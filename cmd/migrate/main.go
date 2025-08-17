package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'")
	}
	direction := os.Args[1]

	// Mengambil string koneksi dari environment variable
	// Sangat disarankan untuk tidak menyimpan kredensial di kode
	pgURI := os.Getenv("PG_URI")
	if pgURI == "" {
		log.Fatal("PG_URI environment variable is not set")
	}

	// Membuat instance migrasi
	// Perhatikan bahwa kita tidak lagi menggunakan sql.Open dan WithInstance
	// go-migrate akan menangani koneksi database internalnya
	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		pgURI,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Memanggil Up() atau Down()
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations applied successfully!")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("Migrations rolled back successfully!")
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}
