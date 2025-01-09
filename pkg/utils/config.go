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
	rawYaml, err := os.ReadFile("sync.yaml")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	err = yaml.Unmarshal([]byte(rawYaml), &_config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func (c *Config) Get() *data.SyncConfig {
	return &_config
}
