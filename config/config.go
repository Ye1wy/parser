package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/lpernett/godotenv"
)

type ConfigProvider interface {
	GetOptHeadless() bool
	GetRequestDelay() string
	GetProxies() []string
	GetCategories() []string
}

type config struct {
	Headless     bool     `env:"HEADLESS" envDefault:"true"`
	RequestDelay string   `env:"REQUEST_DELAY" envDefault:"5"`
	Proxies      []string `env:"PROXIES" envSeparator:","`
	Categories   []string `json:"categories"`
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

	reader, err := os.ReadFile("targets.json")
	if err != nil {
		log.Fatalf("Failed to open catalog of targets: %v,  %s", err, op)
	}

	if err = json.Unmarshal(reader, &cfg); err != nil {
		log.Fatalf("Failed of unmarshal json: %v, %s", err, op)
	}

	return &cfg
}

func (c *config) GetOptHeadless() bool {
	return c.Headless
}

func (c *config) GetRequestDelay() string {
	return c.RequestDelay
}

func (c *config) GetProxies() []string {
	return c.Proxies
}

func (c *config) GetCategories() []string {
	return c.Categories
}
