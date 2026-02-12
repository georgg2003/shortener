package models

import "github.com/google/uuid"

// Структура, которая используется для хранения данных в файловом хранилище.
type ShortURL struct {
	UUID      uuid.UUID `json:"uuid"`
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	IsDeleted bool      `json:"is_deleted"`
}

// Универсальная структура, которая используется для всех операций с URL.
type URLEntity struct {
	CorrelationID string
	ShortURL      string
	OriginalURL   string
	ShortID       string
	IsDeleted     bool
}
