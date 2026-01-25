package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/georgg2003/shortener/internal/repository/audit"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/sirupsen/logrus"
)

type auditFileRepository struct {
	logger   *logrus.Entry
	filename string
	mu       sync.RWMutex
}

func (repo *auditFileRepository) WriteLog(log audit.AuditLog) error {
	data, err := json.Marshal(log)
	if err != nil {
		repo.logger.Error()
		return utils.ErrWrap(err, "failed to marshall an audit log")
	}
	repo.logger.WithFields(logrus.Fields{
		"log":      log,
		"filename": repo.filename,
	}).Info("writing audit log into file")

	repo.mu.Lock()
	defer repo.mu.Unlock()

	f, err := os.OpenFile(repo.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return utils.ErrWrap(err, "failed to open an audit file")
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%s\n", data); err != nil {
		return utils.ErrWrap(err, "failed to write an audit log")
	}

	return nil
}

func New(logger *logrus.Entry, filename string) audit.AuditRepository {
	return &auditFileRepository{
		logger:   logger,
		filename: filename,
		mu:       sync.RWMutex{},
	}
}
