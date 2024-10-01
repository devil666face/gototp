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

func (s *Storage) saveToFile(filename string, data interface{}, encrypt bool) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode data: %w", err)
	}
	if encrypt {
		cryptbytes, err := s.cryptor.Encrypt(buff.Bytes())
		if err != nil {
			return fmt.Errorf("failed encrypt database: %w", err)
		}
		buff.Reset()
		if _, err := buff.Write(cryptbytes); err != nil {
			return fmt.Errorf("failed write bytes: %w", err)
		}
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed save file %s: %w", filename, err)
	}
	return nil
}

func (s *Storage) loadFromFile(filename string, data interface{}, decrypt bool) error {
	rawbytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed read data from file %s: %w", filename, err)
	}
	if decrypt {
		decbytes, err := s.cryptor.Decrypt(rawbytes)
		if err != nil {
			return fmt.Errorf("failed decrypt database: %w", err)
		}
		rawbytes = decbytes
	}
	buff := bytes.NewBuffer(rawbytes)
	dec := gob.NewDecoder(buff)
	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("failed decode data: %w", err)
	}
	return nil
}

func (s *Storage) Save(data *models.Data) error {
	return s.saveToFile(s.filename, data, true)
}

func (s *Storage) Load() (*models.Data, error) {
	var data models.Data
	if err := s.loadFromFile(s.filename, &data, true); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Storage) SaveFile(key *models.Key) error {
	filename := key.Name + ".gototp"
	return s.saveToFile(filename, key, false)
}

func (s *Storage) LoadFile(filename string) (*models.Key, error) {
	var key models.Key
	if err := s.loadFromFile(filename, &key, false); err != nil {
		return nil, err
	}
	return &key, nil
}
