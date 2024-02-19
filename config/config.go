package config

import (
	"errors"
	"os"

	"gemini-chat/utils"
)

type Config struct {
	APIKey string
}

func NewConfig() (*Config, error) {
	var errs []error

	apiKey, err := getEnv("GEMINI_API_KEY", parseString, "")
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}

func getEnv[T any](key string, parser func(value string) (T, error), defaultValue T) (T, error) {
	value, ok := os.LookupEnv(key)
	if ok {
		parsed, err := parser(value)

		return parsed, utils.ErrWrapf(err, "parsing env %s", key)
	}

	return defaultValue, nil
}

func parseString(value string) (string, error) {
	return value, nil
}
