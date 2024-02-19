package utils

import "fmt"

func ErrWrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}

func ErrWrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return ErrWrap(err, fmt.Sprintf(format, args...))
}
