package models

type APIShortenURLRequest struct {
	URL string `json:"url"`
}

type APIShortenURLResponse struct {
	Result string `json:"result"`
}

type APIShortenURLBatchRequestRecord struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type APIShortenURLBatchRequest []APIShortenURLBatchRequestRecord

type APIShortenURLBatchResponseRecord struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type APIShortenURLBatchResponse []APIShortenURLBatchResponseRecord

type APIUserURLsResponseRecord struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type APIUserURLsResponse []APIUserURLsResponseRecord

type APIDeleteUserURLsRequest []string
