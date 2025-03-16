package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Database struct {
	Host         string        `yaml:"DB_HOST" envDefault:"localhost"`
	Port         int           `yaml:"DB_PORT" envDefault:"5432"`
	User         string        `yaml:"DB_USER" envDefault:"postgres"`
	Password     string        `yaml:"DB_PASSWORD" envDefault:"postgres"`
	Database     string        `yaml:"DB_DATABASE" envDefault:"postgres"`
	SSLMode      string        `yaml:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns int           `yaml:"DB_MAX_OPEN_CONNECTS" envDefault:"25"`
	MaxIdleConns int           `yaml:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	MaxTimeLife  time.Duration `yaml:"DB_MAX_TIME_LIFE" envDefault:"24h"`
}

func (d *Database) DBUrl() string {
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s  port=%d sslmode=%s",
		d.User, d.Database, d.Password, d.Host, d.Port, d.SSLMode,
	)
}

type Config struct {
	Env         string   `yaml:"env" env-default:"local"`
	StoragePath string   `yaml:"storage_path" env-required:"true"`
	Database    Database `yaml:"database" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
