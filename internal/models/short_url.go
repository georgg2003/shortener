package models

type ShortURL struct {
	UUID     string `json:"uuid"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}
