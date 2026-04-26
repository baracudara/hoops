package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    Env      string `yaml:"env" env-default:"local"`
    HTTP     `yaml:"http"`
    AuthGRPC `yaml:"auth_grpc"`
}

type HTTP struct {
    Port         int           `yaml:"port" env-default:"8080"`
    Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
    CookieDomain string        `yaml:"cookie_domain" env-default:"localhost"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"720h"`
}

type AuthGRPC struct {
	Host    string        `yaml:"host" env-default:"localhost"`
    Port    int           `yaml:"port" env-default:"44044"`
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