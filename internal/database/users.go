package database

import (
	"context"
	"database/sql"
	"event-app/internal/cache"
	"fmt"
	"time"
)

type UserModel struct {
	DB    *sql.DB
	Cache *cache.Cache
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id"

	err := m.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.Id)
	if err != nil {
		return err
	}

	// Invalidate cache setelah insert
	if m.Cache != nil {
		// Hapus cache untuk events list karena mungkin ada perubahan
		m.Cache.Delete(ctx, "events:list")
	}

	return nil
}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Coba ambil dari cache terlebih dahulu
	if m.Cache != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		var cachedUser User
		err := m.Cache.Get(ctx, cacheKey, &cachedUser)
		if err == nil {
			return &cachedUser, nil
		}
	}

	// Jika tidak ada di cache, ambil dari database
	query := "SELECT id, email, name, password FROM users WHERE id = $1"
	user, err := m.getUser(query, id)
	if err != nil {
		return nil, err
	}

	// Simpan ke cache jika user ditemukan
	if user != nil && m.Cache != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		m.Cache.Set(ctx, cacheKey, user, 30*time.Minute)
	}

	return user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Coba ambil dari cache terlebih dahulu
	if m.Cache != nil {
		cacheKey := "user:email:" + email
		var cachedUser User
		err := m.Cache.Get(ctx, cacheKey, &cachedUser)
		if err == nil {
			return &cachedUser, nil
		}
	}

	// Jika tidak ada di cache, ambil dari database
	query := "SELECT id, email, name, password FROM users WHERE email = $1"
	user, err := m.getUser(query, email)
	if err != nil {
		return nil, err
	}

	// Simpan ke cache jika user ditemukan
	if user != nil && m.Cache != nil {
		cacheKey := "user:email:" + email
		m.Cache.Set(ctx, cacheKey, user, 30*time.Minute)
	}

	return user, nil
}
