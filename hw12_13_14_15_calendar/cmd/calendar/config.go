package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger struct {
		Level string `validate:"required"`
		Path  string `validate:"required"`
	}
	Server struct {
		Host            string        `validate:"required"`
		Port            string        `validate:"required"`
		ShutDownTimeOut time.Duration `validate:"required"`
	}
	DB struct {
		URL string `validate:"required"`
	}
}

func NewConfig(path string) (Config, error) {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("can't read config file: %w", err)
	}

	var config Config

	if err := yaml.Unmarshal(cfg, &config); err != nil {
		return Config{}, fmt.Errorf("can't parse config file: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return Config{}, fmt.Errorf("config invalid: %w", err)
	}

	return config, nil
}
