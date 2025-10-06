package repository

type Repository interface {
	NewUrl(url string, shortUrlId string)
}

type repository struct{}

func New() Repository {
	return repository{}
}
