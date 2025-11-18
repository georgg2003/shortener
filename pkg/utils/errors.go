package utils

import "fmt"

func ErrWrap(err0 error, err1 error) error {
	return fmt.Errorf("%v: %w", err1, err0)
}
