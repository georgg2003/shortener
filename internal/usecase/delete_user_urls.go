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

	for {
		var tasks []models.DeleteUserURLsTask

		select {
		case task := <-uc.deleteUserURLsCh:
			tasks = append(tasks, task)
		case <-ticker.C:
			if len(tasks) == 0 {
				continue
			}

			err := uc.repository.DeleteUserURLs(context.TODO(), tasks)
			if err != nil {
				uc.logger.
					WithFields(logrus.Fields{"tasks_queue": len(tasks)}).
					WithError(err).
					Error("failed to delete user urls")
				continue
			}

			tasks = nil
		}
	}
}
