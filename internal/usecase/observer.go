package usecase

import "time"

type ObserverEvent struct {
	Time        time.Time
	UserID      int64
	OriginalURL string
}

type Observer interface {
	OnNewURL(ObserverEvent)
	OnGetURL(ObserverEvent)
}

type observersByID = map[string]Observer

func (uc *useCase) Observe(id string, observer Observer) (unregister func()) {
	logger := uc.logger.WithField("id", id)
	logger.Info("registered observer")
	if uc.observersByID == nil {
		uc.observersByID = make(observersByID)
	}
	uc.observersByID[id] = observer
	return func() {
		delete(uc.observersByID, id)
		logger.Info("unregistered observer")
	}
}
