package models

import "github.com/google/uuid"

type ShortURL struct {
	UUID     uuid.UUID `json:"uuid"`
	ShortURL string    `json:"short_url"`
	LongURL  string    `json:"long_url"`
}

type URLEntity struct {
	CorrelationID string
	ShortURL      string
	OriginalURL   string
	ShortID       string
}
