package database

import (
	"context"
	"database/sql"
	"event-app/internal/cache"
	"fmt"
	"time"
)

type EventModel struct {
	DB    *sql.DB
	Cache *cache.Cache
}

type Event struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"ownerId"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3"`
}

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := m.DB.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
	if err != nil {
		return err
	}

	// Invalidate cache setelah insert
	if m.Cache != nil {
		// Hapus cache untuk events list dan user cache
		m.Cache.Delete(ctx, "events:list")
		userCacheKey := fmt.Sprintf("user:%d", event.OwnerId)
		m.Cache.Delete(ctx, userCacheKey)
	}

	return nil
}

func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Coba ambil dari cache terlebih dahulu
	if m.Cache != nil {
		var cachedEvents []*Event
		err := m.Cache.Get(ctx, "events:list", &cachedEvents)
		if err == nil {
			return cachedEvents, nil
		}
	}

	// Jika tidak ada di cache, ambil dari database
	query := "SELECT id, owner_id, name, description, date, location FROM events"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event

		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Simpan ke cache
	if m.Cache != nil {
		m.Cache.Set(ctx, "events:list", events, 15*time.Minute)
	}

	return events, nil
}

func (m *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Coba ambil dari cache terlebih dahulu
	if m.Cache != nil {
		cacheKey := fmt.Sprintf("event:%d", id)
		var cachedEvent Event
		err := m.Cache.Get(ctx, cacheKey, &cachedEvent)
		if err == nil {
			return &cachedEvent, nil
		}
	}

	// Jika tidak ada di cache, ambil dari database
	query := "SELECT id, owner_id, name, description, date, location FROM events WHERE id = $1"

	var event Event

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Simpan ke cache jika event ditemukan
	if m.Cache != nil {
		cacheKey := fmt.Sprintf("event:%d", id)
		m.Cache.Set(ctx, cacheKey, &event, 30*time.Minute)
	}

	return &event, nil
}

func (m *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE events SET name = $1, description = $2, date = $3, location = $4 WHERE id = $5"

	_, err := m.DB.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.Id)
	if err != nil {
		return err
	}

	// Invalidate cache setelah update
	if m.Cache != nil {
		// Hapus cache untuk event ini dan events list
		eventCacheKey := fmt.Sprintf("event:%d", event.Id)
		m.Cache.Delete(ctx, eventCacheKey, "events:list")

		// Update cache dengan data baru
		m.Cache.Set(ctx, eventCacheKey, event, 30*time.Minute)
	}

	return nil
}

func (m *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM events WHERE id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Invalidate cache setelah delete
	if m.Cache != nil {
		// Hapus cache untuk event ini dan events list
		eventCacheKey := fmt.Sprintf("event:%d", id)
		m.Cache.Delete(ctx, eventCacheKey, "events:list")
	}

	return nil
}
