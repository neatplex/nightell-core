package config

import (
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/neatplex/nightell-core/internal/http/server/validator"
	"github.com/neatplex/nightell-core/internal/utils"
	"os"
)

const AppName = "Nightell"
const AppVersion = "v0.1.0"

const defaultConfigPath = "configs/main.defaults.json"
const envConfigPath = "configs/main.json"

// Config is the project root level configuration.
type Config struct {
	Development bool   `json:"development"`
	URL         string `json:"url"`
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
	Settings struct {
		AndroidLastVersion          string `json:"android_last_version"`
		AndroidLastSupportedVersion string `json:"android_last_supported_version"`
	} `json:"settings"`
	MySQL struct {
		Host     string `yaml:"host" validate:"required"`
		Port     int    `yaml:"port" validate:"required"`
		Name     string `yaml:"name" validate:"required"`
		User     string `yaml:"user" validate:"required"`
		Password string `yaml:"password" validate:"required"`
		Timeout  int    `yaml:"timeout" validate:"required"`
	} `json:"mysql"`
	S3 struct {
		RoleUsed  bool   `yaml:"role_used"`
		AccessKey string `yaml:"accessKey" validate:"required"`
		SecretKey string `yaml:"secretKey" validate:"required"`
		Region    string `yaml:"region" validate:"required"`
		Bucket    string `yaml:"bucket" validate:"required"`
	} `json:"s3"`
	Mailer struct {
		SmtpServer string `json:"smtp_server"`
		SmtpPort   int    `json:"smtp_port"`
		Username   string `json:"username"`
		Password   string `json:"password"`
	} `json:"mailer"`
	Google struct {
		OAuthClientId string `yaml:"oauth_client_id"`
	} `json:"google"`
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

	fmt.Println("Config:", *c)

	return validator.New().Validate(c)
}

func New() *Config {
	return &Config{}
}
