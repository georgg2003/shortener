package utils

import "fmt"

func ErrWrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
