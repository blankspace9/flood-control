package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Mongo    Mongo         `yaml:"mongo"`
		Window   time.Duration `yaml:"window" env-required:"true"`
		MaxCalls int           `yaml:"max_calls" env-required:"true"`
	}

	Mongo struct {
		Host     string        `env:"MONGODB_HOST"`
		Port     int           `env:"MONGODB_PORT"`
		Username string        `env:"MONGODB_USERNAME"`
		Password string        `env:"MONGODB_PASSWORD"`
		DBName   string        `env:"MONGODB_DBNAME"`
		Timeout  time.Duration `env:"MONGODB_TIMEOUT"`
	}
)

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg.Mongo); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// flag > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = "./config.yaml"
	}

	return res
}
