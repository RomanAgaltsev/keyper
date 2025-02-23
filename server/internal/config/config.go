// Package config provides application configuration building.
package config

import (
    "errors"
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/ilyakaznacheev/cleanenv"
)

// Config contains application configuration.
type Config struct {
    Env      string        `yaml:"env" env-default:"dev"`      // Type of environment to run application - "dev", "prod"
    TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"` // JWT-tokens lifetime

    App   *AppConfig   `yaml:"app"`   // Application configuration
    GRPC  *GRPCConfig  `yaml:"grpc"`  // GRPC-server configuration
    Pprof *PprofConfig `yaml:"pprof"` // Pprof HTTP-server configuration
}

// AppConfig - application configuration.
type AppConfig struct{}

// GRPCConfig - GRPC-server configuration.
type GRPCConfig struct {
    Port    int           `yaml:"port" env-default:":9090"`  // Port of GRPC-server to run
    Timeout time.Duration `yaml:"timeout" env-default:"10h"` // GRPS-server timeout
}

// PprofConfig - Pprof HTTP-server configuration.
type PprofConfig struct {
    Address string `yaml:"address" env-default:":2080"` // Address of pprof HTTP-server to run
}

// Load reads configuration from a file specified with flag or env.
func Load() (*Config, error) {
    var errConfigFile error

    // Get path to the config file from flag/env
    configPath := getConfigPath()
    if configPath == "" {
        errConfigFile = errors.New("config file is not specified")
    }

    // Check config file existence
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        errConfigFile = fmt.Errorf("config file does not exist: %s", configPath)
    }

    if errConfigFile != nil {
        return nil, fmt.Errorf("failed to init config: %w", errConfigFile)
    }

    var cfg Config

    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        return nil, fmt.Errorf("failed to init config: %w", err)
    }

    return &cfg, nil
}

func getConfigPath() string {
    var result string

    flag.StringVar(&result, "config", "", "path to config file")
    flag.Parse()

    if result == "" {
        result = os.Getenv("CONFIG")
    }

    return result
}
