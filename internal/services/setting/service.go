package setting

import (
	"github.com/neatplex/nightell-core/internal/config"
)

type Service struct {
	config *config.Config
}

type Settings struct {
	AndroidLastVersion          string `json:"android_last_version"`
	AndroidLastSupportedVersion string `json:"android_last_supported_version"`
}

func (s *Service) Get() *Settings {
	return &Settings{
		AndroidLastVersion:          s.config.Settings.AndroidLastVersion,
		AndroidLastSupportedVersion: s.config.Settings.AndroidLastSupportedVersion,
	}
}

func New(c *config.Config) *Service {
	return &Service{config: c}
}
