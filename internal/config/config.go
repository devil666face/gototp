package config

import (
	"fmt"
	"gototp/pkg/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	config   = "config.yaml"
	database = "gototp.enc"
)

type Config struct {
	Database string `yaml:"database" env-default:""`
}

func configPrefix() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(u.HomeDir, "AppData", "Roaming", "gototp"), nil
	default:
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return xdgConfigHome, nil
		}
		return filepath.Join(u.HomeDir, ".config", "gototp"), nil
	}
}

func New() (*Config, error) {
	var _config = Config{}

	prefix, err := configPrefix()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(prefix); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(prefix, 0755)
		}
	}
	path := filepath.Join(prefix, config)
	if err := cleanenv.ReadConfig(path, &_config); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config in %s: %w", path, err)
		}
	}
	if _config.Database = fs.UserPath(_config.Database); _config.Database == "" {
		_config.Database = fs.UserPath(filepath.Join(prefix, database))
	}
	return &_config, nil
}
