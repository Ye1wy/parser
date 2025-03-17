package config

import (
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/lpernett/godotenv"
)

type ConfigProvider interface {
	GetProxies() []string
	GetRequestDelay() time.Duration
	GetOptHeadless() bool
}

type config struct {
	Headless     bool          `env:"HEADLESS" envDefault:"true"`
	RequestDelay time.Duration `env:"REQUEST_DELAY" envDefault:"5s"`
	Proxies      []string      `env:"PROXIES" envSeparator:","`
}

func MustLoad() *config {
	op := "config.MustLoad"

	if err := godotenv.Load(); err != nil {
		log.Fatal("op", op, err)
	}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v, op: %s", err, op)
	}

	for i, proxy := range cfg.Proxies {
		cfg.Proxies[i] = strings.TrimSpace(proxy)
	}

	return &cfg
}

func (c *config) GetProxies() []string {
	return c.Proxies
}

func (c *config) GetRequestDelay() time.Duration {
	return c.RequestDelay
}

func (c *config) GetOptHeadless() bool {
	return c.Headless
}
