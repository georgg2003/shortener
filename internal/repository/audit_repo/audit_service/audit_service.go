package audit_service

import (
	"time"

	"github.com/georgg2003/shortener/internal/repository/audit_repo"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type auditServiceRepository struct {
	logger *logrus.Entry
	addr   string
	client *resty.Client
}

func (repo *auditServiceRepository) WriteLog(log audit_repo.AuditLog) error {
	repo.logger.WithFields(logrus.Fields{
		"log":  log,
		"addr": repo.addr,
	}).Info("writing audit log into service")

	_, err := repo.client.R().SetBody(log).Send()
	return err
}

func New(logger *logrus.Entry, addr string) audit_repo.AuditRepository {
	return &auditServiceRepository{
		logger: logger,
		addr:   addr,
		client: resty.New().SetBaseURL(addr).
			SetRetryCount(5).
			SetRetryWaitTime(time.Second),
	}
}
