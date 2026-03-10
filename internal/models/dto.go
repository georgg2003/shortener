package models

// Тело запроса в API ручку для сокращения URL.
//
//	{
//		"url": "https://calendar.mail.ru"
//	}
type APIShortenURLRequest struct {
	// URL для сокращения.
	URL string `json:"url"`
}

// Тело ответа API ручки для сокращения URL.
//
//	{
//		"url": "https://shortener.example.com/d76ascYn"
//	}
type APIShortenURLResponse struct {
	// Сокращенный URL.
	Result string `json:"result"`
}

// Элемент в теле запроса в API ручку для сокращения нескольких URL.
type APIShortenURLBatchRequestRecord struct {
	// ID для сопоставления оригинальных и сокращенных URL между собой.
	CorrelationID string `json:"correlation_id"`
	// URL для сокращения.
	OriginalURL string `json:"original_url"`
}

// Тело запроса в API ручку для сокращения нескольких URL.
//
//	[
//		{
//			"correlation_id": "1",
//			"original_url": "https://calendar.mail.ru"
//		}
//	]
type APIShortenURLBatchRequest []APIShortenURLBatchRequestRecord

// Элемент в теле ответа API ручки для сокращения нескольких URL.
type APIShortenURLBatchResponseRecord struct {
	// ID для сопоставления оригинальных и сокращенных URL между собой.
	CorrelationID string `json:"correlation_id"`
	// Сокращенный URL.
	ShortURL string `json:"short_url"`
}

// Тело ответа API ручки для сокращения нескольких URL.
//
//	[
//		{
//			"correlation_id": "1",
//			"short_id": "https://shortener.example.com/d76ascYn"
//		}
//	]
type APIShortenURLBatchResponse []APIShortenURLBatchResponseRecord

// Элемент в теле ответа API ручки для получения всех URL пользователя.
type APIUserURLsResponseRecord struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// Тело ответа API ручки для получения всех URL пользователя.
//
//	[
//		{
//			"original_url": "https://calendar.mail.ru",
//			"short_url": "https://shortener.example.com/d76ascYn"
//		}
//	]
type APIUserURLsResponse []APIUserURLsResponseRecord

// Тело запроса в API ручку для удаления URL пользователя.
// Список сокращенных ID.
type APIDeleteUserURLsRequest []string

type APIInternalStatsResponse struct {
	URLsCount  int64 `json:"urls"`
	UsersCount int64 `json:"users"`
}
