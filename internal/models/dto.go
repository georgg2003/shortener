package models

type APIShortenURLRequest struct {
	URL string `json:"url"`
}

type APIShortenURLResponse struct {
	Result string `json:"result"`
}
