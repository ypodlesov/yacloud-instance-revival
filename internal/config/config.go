package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GrpcClient    `yaml:"grpc_client"`
	RevivalConfig `yaml:"revival_config"`
}

type GrpcClient struct {
	Address string `yaml:"address"`
}

type Instance struct {
	InstanceId        string        `yaml:"instance_id"`
	CheckHealthPeriod time.Duration `yaml:"check_health_period"`
}

type RevivalConfig struct {
	CheckHealthPeriod time.Duration `yaml:"check_health_period"`
	Instances         []Instance    `yaml:"instances"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	for idx, instance := range cfg.RevivalConfig.Instances {
		if instance.CheckHealthPeriod == 0*time.Second {
			cfg.RevivalConfig.Instances[idx].CheckHealthPeriod = cfg.RevivalConfig.CheckHealthPeriod
		}
	}

	return &cfg
}
