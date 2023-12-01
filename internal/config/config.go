package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"strings"
)

// AppName is the application name.
const AppName = "Nightel"

// AppVersion is the application version.
const AppVersion = "v0.1.0"

// Config is the project root level configuration.
type Config struct {
	Environment Environment `yaml:"environment" validate:"required"`
	Timezone    string      `yaml:"timezone" validate:"required"`
	Logger      Logger      `yaml:"logger" validate:"required"`
	Database    Database    `yaml:"database" validate:"required"`
	HTTPServer  HTTPServer  `yaml:"httpServer" validate:"required"`
	S3          S3          `yaml:"s3" validate:"required"`
}

type Environment string

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
)

// Logger is the logging configuration.
type Logger struct {
	Level string `yaml:"level" validate:"required,lowercase,oneof=debug info warn error fatal panic"`
	Path  string `yaml:"path"`
}

// Database is the database (sql) configuration.
type Database struct {
	Driver   string `yaml:"driver" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Timeout  int    `yaml:"timeout" validate:"required"`
}

// DSN returns the database DSN.
func (d Database) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		"utf8mb4_general_ci",
	)
}

// HTTPServer is http server configuration.
type HTTPServer struct {
	Listen            string `yaml:"listen" validate:"required"`
	ReadTimeout       string `yaml:"readTimeout" validate:"required"`
	WriteTimeout      string `yaml:"writeTimeout" validate:"required"`
	ReadHeaderTimeout string `yaml:"readHeaderTimeout" validate:"required"`
	IdleTimeout       string `yaml:"idleTimeout" validate:"required"`
}

// S3 is file storage by AWS.
type S3 struct {
	AccessKey string `yaml:"accessKey" validate:"required"`
	SecretKey string `yaml:"secretKey" validate:"required"`
	Region    string `yaml:"region" validate:"required"`
	Bucket    string `yaml:"bucket" validate:"required"`
}

// New creates a new configuration instance.
func New(path string) (*Config, error) {
	c := new(Config)

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix(strings.ToLower(AppName))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}

	fmt.Println("configuration", *c)

	return c, validator.New().Struct(c)
}
