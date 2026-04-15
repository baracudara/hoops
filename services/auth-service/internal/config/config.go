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
	Redis    `yaml:"redis"`
	GRPC     `yaml:"grpc"`
	JWT      `yaml:"jwt"`
}



type Postgres struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"root"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"qwerty"`
	DBName   string `yaml:"dbname" env-default:"auth"`
	MinConns int32  `yaml:"min_conns" env-default:"2"`
	MaxConns int32  `yaml:"max_conns" env-default:"10"`
}

type Redis struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"6379"`
	User     string `yaml:"username"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:"qwerty"`
	DB       int    `yaml:"db" env-default:"0"`
}

type JWT struct {
	Secret          string        `yaml:"secret" env:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"240h"`
}

type GRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	cfgPath := fetchConfigPath()

	if cfgPath == "" {
		panic("Config path has not been set")
	}

	return FetchConfig(cfgPath)
}

func FetchConfig(cfgPath string) *Config {
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("No files in the spicified directory" + cfgPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(cfgPath, &config); err != nil {
		panic("Unable to read config of the spicified directory" + err.Error())
	}

	return &config

}



func fetchConfigPath() string {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", ".config/local.yaml", "get config path")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = os.Getenv("CONFIG_PATH")  // ← вот так
	}

	return cfgPath
}

