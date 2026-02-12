// Утилиты для работы с файловым JSON хранилищем
package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/sirupsen/logrus"
)

type SyncMapStorage struct {
	sync.Map
	saveCh chan struct{}
	cfg    *SyncStorageConfig
	logger logrus.FieldLogger
}

type SyncStorageConfig struct {
	FileStoragePath  string
	DebounceDuration *time.Duration
}

func (s *SyncMapStorage) Store(key, value any) {
	s.Map.Store(key, value)
	select {
	case s.saveCh <- struct{}{}:
	default:
	}
}

func (s *SyncMapStorage) recoverFromFile() {
	file, err := os.Open(s.cfg.FileStoragePath)
	if os.IsNotExist(err) {
		s.logger.WithError(err).Info("file does not exist")
		return
	}
	if err != nil {
		s.logger.WithError(err).Error("failed to open file storage")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	tmp := make([]models.ShortURL, 0)
	if err := decoder.Decode(&tmp); err != nil {
		s.logger.WithError(err).Error("failed to decode data")
		return
	}

	for _, k := range tmp {
		s.Store(k.ShortURL, k)
	}

	s.logger.Infof("successfully recovered data from file, rows: %d", len(tmp))
}

func (s *SyncMapStorage) syncWorker() {
	for range s.saveCh {
		if s.cfg.DebounceDuration != nil {
			duration := *s.cfg.DebounceDuration
			time.Sleep(duration)
		}

		tmp := make([]models.ShortURL, 0)

		s.Map.Range(func(key, value any) bool {
			v, ok := value.(models.ShortURL)
			if ok {
				tmp = append(tmp, v)
			}
			return true
		})

		bytes, err := json.MarshalIndent(tmp, "", "  ")
		if err != nil {
			s.logger.WithError(err).Error("failed marshal data")
			continue
		}

		err = os.WriteFile(s.cfg.FileStoragePath, bytes, 0666)
		if err != nil {
			s.logger.WithError(err).Error("failed to write data to file")
		}

		s.logger.WithFields(
			logrus.Fields{
				"file_path": s.cfg.FileStoragePath,
				"bytes":     len(bytes),
			},
		).Infof("successfully written data to file storage")
	}
}

func New(
	cfg *SyncStorageConfig,
	logger logrus.FieldLogger,
) *SyncMapStorage {
	saveCh := make(chan struct{}, 1)

	store := SyncMapStorage{
		saveCh: saveCh,
		cfg:    cfg,
		logger: logger,
	}
	store.recoverFromFile()
	go store.syncWorker()

	return &store
}
