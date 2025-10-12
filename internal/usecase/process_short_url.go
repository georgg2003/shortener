package usecase

func (uc useCase) ProcessShortURL(id string) (string, error) {
	longURL, err := uc.repository.GetLongURL(id)
	return longURL, err
}
