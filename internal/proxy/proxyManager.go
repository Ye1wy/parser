package proxy

import (
	"log/slog"
	"math/rand/v2"
	"parser/config"
)

type proxyManager struct {
	log     *slog.Logger
	proxies []string
}

func NewProxyManager(cfg config.ConfigProvider, logger *slog.Logger) *proxyManager {
	return &proxyManager{
		log:     logger,
		proxies: cfg.GetProxies(),
	}
}

func (pm *proxyManager) GetRandomProxy() string {
	op := "proxy.proxyManager.GetRandomProxy"

	if len(pm.proxies) == 0 {
		pm.log.Warn("No proxies", "op", op)
		return ""
	}

	return pm.proxies[rand.IntN(len(pm.proxies))]
}
