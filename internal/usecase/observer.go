package usecase

import "time"

const NewURLAction = "shorten"
const FollowURLAction = "follow"

type ObserverEvent struct {
	Time        time.Time
	UserID      int64
	OriginalURL string
	Action      string
}

type Observer interface {
	OnNewURL(ObserverEvent)
	OnGetURL(ObserverEvent)
}

type ObserversByID = map[string]Observer

func (uc *useCase) Observe(id string, observer Observer) (unregister func()) {
	if uc.observersByID == nil {
		uc.observersByID = make(ObserversByID)
	}
	uc.observersByID[id] = observer
	return func() {
		delete(uc.observersByID, id)
	}
}
