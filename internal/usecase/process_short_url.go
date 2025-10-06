package usecase

func (s useCase) ProcessShortURL(id string) (string, error) {
	longURL, err := s.repository.GetLongURL(id)
	return longURL, err
}
