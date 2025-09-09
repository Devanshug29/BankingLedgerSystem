package config

import (
	flags "BankingLedgerSystem/internal/flags"
	"context"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var cfg *Config

// InitConfigs loads configuration once based on deployment mode.
func InitConfigs(ctx context.Context, fileNames ...string) {
	deploymentMode := flags.GetDeploymentMode()
	base := deploymentMode.GetConfigPath()

	var c Config
	for _, f := range fileNames {
		path := fmt.Sprintf("%s/%s", base, f)

		file, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read config file (%s): %v", path, err)
		}

		if err := yaml.Unmarshal(file, &c); err != nil {
			log.Fatalf("failed to parse config file (%s): %v", path, err)
		}
	}

	cfg = &c
}

func GetConfig() *Config {
	if cfg == nil {
		panic("config not initialised, call InitConfigs() first")
	}
	return cfg
}
