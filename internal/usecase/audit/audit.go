package audit

import (
	"slices"
	"strconv"

	audit_repo "github.com/georgg2003/shortener/internal/repository/audit"
	"github.com/georgg2003/shortener/internal/usecase"
	"github.com/sirupsen/logrus"
)

const NewURLAction = "shorten"
const FollowURLAction = "follow"

type AuditObserver struct {
	repositories []audit_repo.AuditRepository
	logger       *logrus.Entry
}

func convertEventToLog(ev usecase.ObserverEvent, action string) audit_repo.AuditLog {
	return audit_repo.AuditLog{
		Timestamp: ev.Time.Unix(),
		Action:    action,
		UserID:    strconv.FormatInt(ev.UserID, 10),
		URL:       ev.OriginalURL,
	}
}

func (obs *AuditObserver) OnNewURL(ev usecase.ObserverEvent) {
	for repo := range slices.Values(obs.repositories) {
		go func() {
			if err := repo.WriteLog(convertEventToLog(ev, NewURLAction)); err != nil {
				obs.logger.WithError(err).Error("failed to write new url log")
			}
		}()
	}
}

func (obs *AuditObserver) OnGetURL(ev usecase.ObserverEvent) {
	for repo := range slices.Values(obs.repositories) {
		go func() {
			if err := repo.WriteLog(convertEventToLog(ev, FollowURLAction)); err != nil {
				obs.logger.WithError(err).Error("failed to write follow url log")
			}
		}()
	}
}

func NewAuditObserver(repos []audit_repo.AuditRepository, logger *logrus.Entry) usecase.Observer {
	return &AuditObserver{
		repositories: repos,
		logger:       logger,
	}
}
