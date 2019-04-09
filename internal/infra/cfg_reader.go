package infra

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/suzuki-shunsuke/akoi/internal/domain"
)

type (
	// ConfigReader implements domain.ConfigReader .
	ConfigReader struct{}
)

// Read implements domain.ConfigReader .
func (reader ConfigReader) Read(dst string) (domain.Config, error) {
	cfg := domain.Config{}
	f, err := os.Open(dst)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&cfg)
	return cfg, err
}
