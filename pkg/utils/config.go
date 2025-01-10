package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/openshift-splat-team/ci-secret-sync/data"
	"gopkg.in/yaml.v2"
)

type Config struct {
}

var (
	_config data.SyncConfig
)

func init() {

}

func LoadConfig(path string) error {
	rawYaml, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return fmt.Errorf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal([]byte(rawYaml), &_config)
	if err != nil {
		fmt.Println("error unmarshalling config file:", err)
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}
	return nil
}

func (c *Config) Get() *data.SyncConfig {
	return &_config
}
