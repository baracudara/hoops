package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    Env      string `yaml:"env" env-default:"local"`
    Postgres `yaml:"postgres"`
    GRPC     `yaml:"grpc"`
}

type Postgres struct {
    Host     string `yaml:"host" env-default:"localhost"`
    Port     int    `yaml:"port" env-default:"5432"`
    User     string `yaml:"user" env-default:"root"`
    Password string `yaml:"password" env-default:"qwerty"`
    DBName   string `yaml:"dbname" env-default:"player"`
    MinConns int32  `yaml:"min_conns" env-default:"2"`
    MaxConns int32  `yaml:"max_conns" env-default:"10"`
}

type GRPC struct {
    Port    int           `yaml:"port" env-default:"44045"`
    Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
    cfgPath := fetchConfigPath()
    if cfgPath == "" {
        panic("config path not set")
    }
    return fetchConfig(cfgPath)
}

func fetchConfig(cfgPath string) *Config {
    if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
        panic("config file not found: " + cfgPath)
    }

    var cfg Config
    if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
        panic("failed to read config: " + err.Error())
    }

    return &cfg
}

func fetchConfigPath() string {
    var cfgPath string
    flag.StringVar(&cfgPath, "config", "config/local.yaml", "path to config")
    flag.Parse()

    if cfgPath == "" {
        cfgPath = os.Getenv("CONFIG_PATH")
    }

    return cfgPath
}