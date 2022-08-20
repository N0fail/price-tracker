//go:build integration
// +build integration

package config

import (
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "QA"

type Config struct {
	Host       string `split_words:"true" default:":8081"`
	DbHost     string `split_words:"true" default:"localhost"`
	DbPort     int    `split_words:"true" default:"5432"`
	DbUser     string `split_words:"true" default:"test"`
	DbPassword string `split_words:"true" default:"test"`
	DbName     string `split_words:"true" default:"test"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	return cfg, err
}
