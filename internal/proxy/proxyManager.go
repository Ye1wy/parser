package proxy

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"parser/internal/logger"
)

const (
	proxyFileName     = "./config/proxy_list.json"
	ExtensionFilePath = "./config/proxy.zip"
)

type ProxyManager struct {
	logger  *slog.Logger
	proxies []proxyes
}

type proxyes struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProtocolType string `json:"type"`
}

func NewProxyManager(logger *slog.Logger) *ProxyManager {
	return &ProxyManager{
		logger: logger,
	}
}

func (pm *ProxyManager) LoadProxy() error {
	op := "proxy.proxyManager.LoadProxy"

	file, err := os.Open(proxyFileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			pm.logger.Error("Failed to close file", logger.Err(err), "op", op)
			return
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&pm.proxies)
	if err != nil {
		pm.logger.Error("Failed to close file", logger.Err(err), "op", op)
		return err
	}

	return nil
}

func (pm *ProxyManager) GetRandomProxy() *proxyes {
	if len(pm.proxies) == 0 {
		return nil
	}

	randomNum := rand.IntN(len(pm.proxies))

	return &pm.proxies[randomNum]
}

func (pm *ProxyManager) SetProxyExt(proxyType string, proxyHost string, proxyPort int, proxyUser string, proxyPass string) error {
	op := "proxy.proxyManager.SetProxyExt"

	manifestJson := `
	{
		"version": "1.0.0",
		"manifest_version": 2,
		"name": "Chrome Proxy",
		"permissions": [
			"proxy",
			"tabs",
			"unlimitedStorage",
			"storage",
			"<all_urls>",
			"webRequest",
			"webRequestBlocking"
		],
		"background": {
			"scripts": ["background.js"]
		},
		"minimum_chrome_version":"22.0.0"
	}
	`
	backgroundJs := fmt.Sprintf(`
	var config = {
			mode: "fixed_servers",
			rules: {
			singleProxy: {
				scheme: "%s",
				host: "%s",
				port: parseInt(%d)
			},
			bypassList: ["localhost"]
			}
		};
		
	function callbackFn(details) {
		return {
			authCredentials: {
				username: "%s",
				password: "%s"
			}
		};
	}
	
	chrome.webRequest.onAuthRequired.addListener(
				callbackFn,
				{urls: ["<all_urls>"]},
				['blocking']
	);
	`, proxyType, proxyHost, proxyPort, proxyUser, proxyPass)

	zipFile, err := os.Create(ExtensionFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err := writeFileToZip(zipWriter, "manifest.json", manifestJson); err != nil {
		pm.logger.Error("Failed in write file \"manifest.json\" to zip", logger.Err(err), "op", op)
		if err := os.Remove(ExtensionFilePath); err != nil {
			pm.logger.Error("Failed remove file \"proxt.zip\" to zip", logger.Err(err), "op", op)
		}

		return err
	}

	if err := writeFileToZip(zipWriter, "background.js", backgroundJs); err != nil {
		pm.logger.Error("Failed in write file \"background.js\" to zip", logger.Err(err), "op", op)
		if err := os.Remove(ExtensionFilePath); err != nil {
			pm.logger.Error("Failed remove file \"proxt.zip\" to zip", logger.Err(err), "op", op)
		}

		return err
	}

	return nil
}

func writeFileToZip(zipWriter *zip.Writer, filename, content string) error {
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
