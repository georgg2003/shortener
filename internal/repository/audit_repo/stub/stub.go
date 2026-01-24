package audit_stub

import (
	"encoding/json"

	"github.com/georgg2003/shortener/internal/repository/audit_repo"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/sirupsen/logrus"
)

type auditStub struct {
	logger *logrus.Entry
}

func (s *auditStub) WriteLog(log audit_repo.AuditLog) error {
	_, err := json.Marshal(log)
	if err != nil {
		s.logger.Error()
		return utils.ErrWrap(err, "failed to marshall audit log")
	}
	s.logger.WithField("log", log).Info("audit WriteLog called")
	return nil
}

func New(logger *logrus.Entry) *auditStub {
	logger.Warn("using audit stub")
	return &auditStub{logger: logger}
}
