package config

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/lpernett/godotenv"
)

type Config struct {
	Headless     string   `env:"HEADLESS" envDefault:"true"`
	RequestDelay string   `env:"REQUEST_DELAY" envDefault:"5"`
	Proxies      []string `env:"PROXIES" envSeparator:","`
}

func MustLoad() Config {
	op := "config.MustLoad"

	if err := godotenv.Load(); err != nil {
		log.Fatal("op", op, err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v, op: %s", err, op)
	}

	for i, proxy := range cfg.Proxies {
		cfg.Proxies[i] = strings.TrimSpace(proxy)
	}

	return cfg
}
