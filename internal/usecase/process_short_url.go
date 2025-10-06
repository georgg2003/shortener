package usecase

func (s useCase) ProcessShortURL(id string) (string, error) {
	longUrl, err := s.repository.GetLongUrl(id)
	return longUrl, err
}
