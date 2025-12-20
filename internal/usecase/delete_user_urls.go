package usecase

import (
	"context"
	"time"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/sirupsen/logrus"
)

func (uc *useCase) DeleteUserURLs(ctx context.Context, shortIDs []string) error {
	userID, ok := contextlib.GetUserID(ctx)
	if !ok {
		return errUserNotFound
	}

	uc.deleteUserURLsCh <- models.DeleteUserURLsTask{
		UserID:   userID,
		ShortIDs: shortIDs,
	}

	return nil
}

func (uc *useCase) deleteUserURLsWorker() {
	ticker := time.NewTicker(10 * time.Second)
	var tasks []models.DeleteUserURLsTask

	for {
		select {
		case task := <-uc.deleteUserURLsCh:
			tasks = append(tasks, task)
			uc.logger.
				WithFields(logrus.Fields{"tasks_queue": len(tasks)}).
				Info("added a new task to queue")
		case <-ticker.C:
			logger := uc.logger.WithFields(logrus.Fields{"tasks_queue": len(tasks)})
			logger.Info("started delete user urls operation")
			if len(tasks) == 0 {
				continue
			}

			err := uc.repository.DeleteUserURLs(context.Background(), tasks)
			if err != nil {
				logger.WithError(err).Error("failed to delete user urls")
				continue
			}
			logger.Info("successfully deleted user urls")

			tasks = nil
		}
	}
}
