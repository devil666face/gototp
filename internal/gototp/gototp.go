package gototp

import (
	"crypto/sha256"
	"fmt"
	"gototp/internal/config"
	"gototp/internal/crypt"
	"gototp/internal/database"
	"gototp/internal/models"
	"io"
)

type Gototp struct {
	Data    *models.Data
	config  *config.Config
	storage *database.Storage
}

func GenHash(input string) []byte {
	h := sha256.New()
	io.WriteString(h, input)
	return h.Sum(nil)
}

func New(_passphrase string) (*Gototp, error) {
	cryptor, err := crypt.New(GenHash(_passphrase))
	if err != nil {
		return nil, fmt.Errorf("failed to init key: %w", err)
	}
	_config, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	_storage, err := database.New(
		_config.Database,
		cryptor,
	)
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

func (g *Gototp) SaveFile(key *models.Key) error {
	return g.storage.SaveFile(key)
}

func (g *Gototp) Load() error {
	_data, err := g.storage.Load()
	if err != nil {
		return err
	}
	g.Data = _data
	return nil
}

func (g *Gototp) LoadFile(filename string) (*models.Key, error) {
	return g.storage.LoadFile(filename)
}
