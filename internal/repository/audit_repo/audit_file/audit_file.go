package audit_file

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/georgg2003/shortener/internal/repository/audit_repo"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/sirupsen/logrus"
)

type auditFileRepository struct {
	logger   *logrus.Entry
	filename string
	mu       sync.RWMutex
}

func (repo *auditFileRepository) WriteLog(log audit_repo.AuditLog) error {
	data, err := json.Marshal(log)
	if err != nil {
		repo.logger.Error()
		return utils.ErrWrap(err, "failed to marshall audit log")
	}
	repo.logger.WithFields(logrus.Fields{
		"log":      log,
		"filename": repo.filename,
	}).Info("writing audit log into file")

	repo.mu.Lock()
	defer repo.mu.Unlock()
	return os.WriteFile(repo.filename, data, os.ModeAppend)
}

func New(logger *logrus.Entry, filename string) audit_repo.AuditRepository {
	return &auditFileRepository{
		logger:   logger,
		filename: filename,
		mu:       sync.RWMutex{},
	}
}
