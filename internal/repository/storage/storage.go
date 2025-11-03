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

func (r *SyncMapStorage) recoverFromFile() {
	file, err := os.Open(r.cfg.FileStoragePath)
	if err != nil {
		r.logger.WithError(err).Error("failed to open file storage")
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	tmp := make([]models.ShortURL, 0)
	if err := decoder.Decode(&tmp); err != nil {
		r.logger.WithError(err).Error("failed to decode data")
		return
	}

	for _, k := range tmp {
		r.Store(k.ShortURL, k)
	}

	r.logger.Infof("successfully recovered data from file, rows: %d", len(tmp))
}

func (r *SyncMapStorage) syncWorker() {
	for range r.saveCh {
		duration := 300 * time.Millisecond
		if r.cfg.DebounceDuration != nil {
			duration = *r.cfg.DebounceDuration
		}

		time.Sleep(duration)
		tmp := make([]models.ShortURL, 0)

		r.Map.Range(func(key, value any) bool {
			v, ok := value.(models.ShortURL)
			if ok {
				tmp = append(tmp, v)
			}
			return true
		})

		bytes, err := json.MarshalIndent(tmp, "", "  ")
		if err != nil {
			r.logger.WithError(err).Error("failed marshal data")
			continue
		}

		err = os.WriteFile(r.cfg.FileStoragePath, bytes, 0644)
		if err != nil {
			r.logger.WithError(err).Error("failed to write data to file")
		}
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
