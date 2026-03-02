package utils

import (
	"io"

	"github.com/sirupsen/logrus"
)

func SafeClose(c io.Closer, l logrus.FieldLogger) {
	if err := c.Close(); err != nil {
		l.WithError(err).Error("failed to close")
	}
}
