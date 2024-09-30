package gototp

import (
	"fmt"
	"gototp/internal/config"
	"gototp/internal/database"
	"gototp/internal/models"
)

type Gototp struct {
	Data    *models.Data
	config  *config.Config
	storage *database.Storage
}

func New() (*Gototp, error) {
	_config, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	_storage, err := database.New(_config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to load database: %w", err)
	}
	_data, err := _storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to parse database: %w", err)
	}
	return &Gototp{
		Data:    _data,
		config:  _config,
		storage: _storage,
	}, nil
}

func (g *Gototp) Save() error {
	return g.storage.Save(g.Data)
}

func (g *Gototp) Load() error {
	_data, err := g.storage.Load()
	if err != nil {
		return err
	}
	g.Data = _data
	return nil
}
