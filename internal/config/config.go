package config

import (
	"encoding/json"
	"github.com/neatplex/nightel-core/internal/utils"
	"github.com/pkg/errors"
	"os"
)

const AppName = "Nightel"
const AppVersion = "v0.1.0"

const defaultConfigPath = "configs/main.defaults.json"
const envConfigPath = "configs/main.json"

// Config is the project root level configuration.
type Config struct {
	Development bool `json:"development"`
	Logger      struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	} `json:"logger"`
	HttpServer struct {
		Host              string `json:"host"`
		Port              int    `json:"port"`
		ReadTimeout       int    `json:"read_timeout"`
		WriteTimeout      int    `json:"write_timeout"`
		ReadHeaderTimeout int    `json:"read_header_timeout"`
		IdleTimeout       int    `json:"idle_timeout"`
	} `json:"http_server"`
	HttpClient struct {
		Timeout int `json:"timeout"`
	} `json:"http_client"`
	Database struct {
		Driver   string `yaml:"driver" validate:"required"`
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required"`
		Name     string `yaml:"name" validate:"required"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Timeout  int    `yaml:"timeout" validate:"required"`
	} `json:"database"`
	S3 struct {
		AccessKey string `yaml:"accessKey" validate:"required"`
		SecretKey string `yaml:"secretKey" validate:"required"`
		Region    string `yaml:"region" validate:"required"`
		Bucket    string `yaml:"bucket" validate:"required"`
	} `json:"s3"`
}

func (c *Config) Init() error {
	content, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return errors.WithStack(err)
	}

	if utils.FileExist(envConfigPath) {
		content, err = os.ReadFile(envConfigPath)
		if err != nil {
			return errors.WithStack(err)
		}
		err = json.Unmarshal(content, &c)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func New() *Config {
	return &Config{}
}
