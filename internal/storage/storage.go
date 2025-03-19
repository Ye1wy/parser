package storage

import (
	"encoding/json"
	"log/slog"
	"os"
	"parser/config"
	"parser/internal/logger"
	"parser/internal/models"
)

type Storage interface {
	CreateFile(category string) (string, error)
	ReadFile(path string) (*os.File, error)
	Save(products []models.Product)
}

type storageJson struct {
	cfg config.ConfigProvider
	log *slog.Logger
}

func NewStorageJson(config config.ConfigProvider, logger *slog.Logger) *storageJson {
	return &storageJson{
		cfg: config,
		log: logger,
	}
}

func (sJ *storageJson) CreateFile(category string) (string, error) {
	op := "storage.storageJson.CreateFile"

	if _, err := os.Stat(sJ.cfg.GetPathOutputData()); os.IsNotExist(err) {
		os.Mkdir(sJ.cfg.GetPathOutputData(), os.ModePerm)
		sJ.log.Info("Created", "dir", sJ.cfg.GetPathOutputData(), "op", op)
	}

	fileName := sJ.cfg.GetPathOutputData() + category + ".json"

	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		_, err = os.Create(fileName)
		if err != nil {
			sJ.log.Error("Cant create file", "file name:", fileName, logger.Err(err), "op", op)
			return "", err
		}

		sJ.log.Info("Created", "file", fileName, "op", op)
	}

	sJ.log.Info("Path is saved", "file name", fileName, "op", op)

	return fileName, nil
}

func (sJ *storageJson) ReadFile(path string) (*os.File, error) {
	op := "storage.storageJson.ReadFile"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		sJ.log.Error("Path is not exist", logger.Err(err), "op", op)
		return nil, err
	}

	var file *os.File
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		sJ.log.Error("Cant open file", "file name:", path, logger.Err(err), "op", op)
		return nil, err
	}

	sJ.log.Info("Opened", "file", path, "op", op)

	return file, nil
}

func (sJ *storageJson) ClearFile(path string) error {
	op := "storage.storageJson.ClearFile"

	if err := os.Truncate(path, 0); err != nil {
		sJ.log.Error("Failed to truncate", logger.Err(err), "op", op)
		return err
	}

	sJ.log.Info("Cleared", "file", path, "op", op)
	return nil
}

func (sJ *storageJson) Save(product []models.Product, file *os.File) {
	op := "storage.storageJson.Save"

	encoder := json.NewEncoder(file)
	encoder.SetIndent(" ", " ")
	if err := encoder.Encode(product); err != nil {
		sJ.log.Error("Failed encode product to json file", logger.Err(err), "op", op)
		return
	}

	sJ.log.Info("All products is saved in file", "op", op)
}
