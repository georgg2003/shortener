package audit

import (
	"strconv"

	"github.com/georgg2003/shortener/internal/repository/audit_repo"
	"github.com/georgg2003/shortener/internal/usecase"
)

const NewURLAction = "shorten"
const FollowURLAction = "follow"

type AuditObserver struct {
	repository audit_repo.AuditRepository
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
	obs.repository.WriteLog(convertEventToLog(ev, NewURLAction))
}

func (obs *AuditObserver) OnGetURL(ev usecase.ObserverEvent) {
	obs.repository.WriteLog(convertEventToLog(ev, FollowURLAction))
}

func NewAuditObserver(repo audit_repo.AuditRepository) usecase.Observer {
	return &AuditObserver{
		repository: repo,
	}
}
