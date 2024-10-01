package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"gototp/internal/crypt"
	"gototp/internal/models"
)

type Storage struct {
	filename string
	cryptor  *crypt.Sync
}

func New(
	_filename string,
	_cryptor *crypt.Sync,
) (*Storage, error) {
	var data = models.Data{}
	_storage := &Storage{
		filename: _filename,
		cryptor:  _cryptor,
	}
	if _, err := os.Stat(_storage.filename); os.IsNotExist(err) {
		if err := _storage.Save(&data); err != nil {
			return nil, fmt.Errorf("create database: %w", err)
		}

	}
	return _storage, nil
}

func (s *Storage) saveToFile(filename string, data interface{}) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode data: %w", err)
	}
	cryptbytes, err := s.cryptor.Encrypt(buff.Bytes())
	if err != nil {
		return fmt.Errorf("failed encrypt database: %w", err)
	}
	if err := os.WriteFile(filename, cryptbytes, 0644); err != nil {
		return fmt.Errorf("failed save file %s: %w", filename, err)
	}
	return nil
}

func (s *Storage) Save(data *models.Data) error {
	return s.saveToFile(s.filename, data)
}

func (s *Storage) loadFromFile(filename string, data interface{}) error {
	rawbytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed read data from file %s: %w", filename, err)
	}
	decbytes, err := s.cryptor.Decrypt(rawbytes)
	if err != nil {
		return fmt.Errorf("failed decrypt database: %w", err)
	}
	buff := bytes.NewBuffer(decbytes)
	dec := gob.NewDecoder(buff)
	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("failed decode data: %w", err)
	}
	return nil
}

func (s *Storage) Load() (*models.Data, error) {
	var data models.Data
	if err := s.loadFromFile(s.filename, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
