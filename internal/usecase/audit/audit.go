package audit

import (
	"slices"
	"strconv"

	audit_repo "github.com/georgg2003/shortener/internal/repository/audit"
	"github.com/georgg2003/shortener/internal/usecase"
)

const NewURLAction = "shorten"
const FollowURLAction = "follow"

type AuditObserver struct {
	repositories []audit_repo.AuditRepository
}

func convertEventToLog(ev usecase.ObserverEvent, action string) audit_repo.AuditLog {
	return audit_repo.AuditLog{
		Timestamp: ev.Time.Second(),
		Action:    action,
		UserID:    strconv.FormatInt(ev.UserID, 10),
		URL:       ev.OriginalURL,
	}
}

func (obs *AuditObserver) OnNewURL(ev usecase.ObserverEvent) {
	if obs.repositories == nil {
		return
	}
	for repo := range slices.Values(obs.repositories) {
		repo.WriteLog(convertEventToLog(ev, NewURLAction))
	}
}

func (obs *AuditObserver) OnGetURL(ev usecase.ObserverEvent) {
	if obs.repositories == nil {
		return
	}
	for repo := range slices.Values(obs.repositories) {
		repo.WriteLog(convertEventToLog(ev, FollowURLAction))
	}
}

func NewAuditObserver(repos []audit_repo.AuditRepository) usecase.Observer {
	return &AuditObserver{
		repositories: repos,
	}
}
