package database

import (
	"database/sql"
	"event-app/internal/cache"
	"event-app/internal/env"
	"log"
)

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeeModel
	Cache     *cache.Cache
}

func NewModels(db *sql.DB) Models {
	// Inisialisasi Redis cache
	redisConfig := env.GetRedisConfig()
	redisCache, err := cache.New(redisConfig.GetRedisAddr(), redisConfig.Password, redisConfig.DB)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without cache.", err)
		redisCache = nil
	} else {
		log.Println("Successfully connected to Redis cache")
	}

	return Models{
		Users:     UserModel{DB: db, Cache: redisCache},
		Events:    EventModel{DB: db, Cache: redisCache},
		Attendees: AttendeeModel{DB: db},
		Cache:     redisCache,
	}
}
