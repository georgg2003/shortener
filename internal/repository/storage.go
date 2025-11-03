package repository

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/georgg2003/shortener/internal/config"
	"github.com/georgg2003/shortener/internal/models"
	"github.com/sirupsen/logrus"
)

type storage struct {
	sync.Map
	saveCh chan struct{}
	cfg    *config.Config
	logger *logrus.Logger
}

func (s *storage) Store(key, value any) {
	s.Map.Store(key, value)
	select {
	case s.saveCh <- struct{}{}:
	default:
	}
}

func (r *storage) recoverFromFile() {
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

func (r *storage) syncWorker() {
	for range r.saveCh {
		time.Sleep(300 * time.Millisecond)
		tmp := make([]models.ShortURL, 0)

		r.Map.Range(func(key, value any) bool {
			tmp = append(tmp, value.(models.ShortURL))
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

func NewStorage(
	cfg *config.Config,
	logger *logrus.Logger,
) *storage {
	saveCh := make(chan struct{}, 1)

	store := storage{
		saveCh: saveCh,
		cfg:    cfg,
		logger: logger,
	}
	store.recoverFromFile()
	go store.syncWorker()

	return &store
}
