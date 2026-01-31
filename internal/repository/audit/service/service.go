package service

import (
	"slices"
	"time"

	"github.com/georgg2003/shortener/internal/repository/audit"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type auditServiceRepository struct {
	logger *logrus.Entry
	addr   string
	client *resty.Client
}

func (repo *auditServiceRepository) WriteLog(log audit.AuditLog) error {
	repo.logger.WithFields(logrus.Fields{
		"log":  log,
		"addr": repo.addr,
	}).Info("writing audit log into service")

	_, err := repo.client.R().SetBody(log).Post("")
	return err
}

type ServiceOption func(*auditServiceRepository)

func WithRetryCount(n int) ServiceOption {
	return func(asr *auditServiceRepository) {
		asr.client.SetRetryCount(n)
	}
}

func WithRetryWaitTime(wt time.Duration) ServiceOption {
	return func(asr *auditServiceRepository) {
		asr.client.SetRetryWaitTime(wt)
	}
}

func New(logger *logrus.Entry, addr string, opts ...ServiceOption) audit.AuditRepository {
	repo := &auditServiceRepository{
		logger: logger,
		addr:   addr,
		client: resty.New().
			SetLogger(logger).
			SetBaseURL(addr),
	}
	for opt := range slices.Values(opts) {
		opt(repo)
	}
	return repo
}
