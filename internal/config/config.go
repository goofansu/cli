package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	ServiceMiniflux = "miniflux"
	ServiceLinkding = "linkding"
	ServiceWallabag = "wallabag"
)

type ServiceConfig struct {
	Endpoint string `toml:"endpoint"`
	APIKey   string `toml:"api_key"`
}

type WallabagConfig struct {
	Endpoint     string `toml:"endpoint"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
}

type Config struct {
	Miniflux ServiceConfig  `toml:"miniflux"`
	Linkding ServiceConfig  `toml:"linkding"`
	Wallabag WallabagConfig `toml:"wallabag"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "mlwcli", "auth.toml"), nil
}

func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	return &cfg, err
}

func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func RemoveService(service string) error {
	cfg, err := Load()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	switch service {
	case ServiceMiniflux:
		cfg.Miniflux = ServiceConfig{}
	case ServiceLinkding:
		cfg.Linkding = ServiceConfig{}
	case ServiceWallabag:
		cfg.Wallabag = WallabagConfig{}
	default:
		return nil
	}

	if err := Save(cfg); err != nil {
		return err
	}

	if cfg.Miniflux.Endpoint == "" && cfg.Linkding.Endpoint == "" && cfg.Wallabag.Endpoint == "" {
		path, err := GetConfigPath()
		if err != nil {
			return err
		}
		return os.Remove(path)
	}

	return nil
}
