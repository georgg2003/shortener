package stub

import (
	"github.com/georgg2003/shortener/internal/repository/audit"
	"github.com/sirupsen/logrus"
)

type auditStub struct {
	logger *logrus.Entry
}

func (s *auditStub) WriteLog(log audit.AuditLog) error {
	s.logger.WithField("log", log).Info("audit WriteLog called")
	return nil
}

func New(logger *logrus.Entry) audit.AuditRepository {
	logger.Warn("using audit stub")
	return &auditStub{logger: logger}
}
