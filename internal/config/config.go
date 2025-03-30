package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Host         string        `env-default:"localhost" yaml:"dbHost"`
	Port         int           `env-default:"5432"      yaml:"dbPort"`
	User         string        `env-default:"postgres"  yaml:"dbUser"`
	Password     string        `env-default:"postgres"  yaml:"dbPassword"`
	Database     string        `env-default:"postgres"  yaml:"dbDatabase"`
	SSLMode      string        `env-default:"disable"   yaml:"dbSslmode"`
	MaxOpenConns int           `env-default:"25"        yaml:"dbMaxOpenConnects"`
	MaxIdleConns int           `env-default:"25"        yaml:"dbMaxIdleConns"`
	MaxTimeLife  time.Duration `env-default:"24h"       yaml:"dbMaxTimeLife"`
}

func (d *Database) DBUrl() string {
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s  port=%d sslmode=%s",
		d.User, d.Database, d.Password, d.Host, d.Port, d.SSLMode,
	)
}

type Auth struct {
	JwtSecret string        `env-default:"secret" yaml:"jwtSecret"`
	ExpTime   time.Duration `env-default:"24h"    yaml:"expTime"`
}

type Application struct {
	Port int `env-default:"8000" yaml:"port"`
}
type Config struct {
	Env         string      `env-default:"local" yaml:"env"`
	StoragePath string      `env-default:"true"  yaml:"storage_path"`
	Database    Database    `env-default:"true"  yaml:"database"`
	Auth        Auth        `env-default:"true"  yaml:"auth"`
	Application Application `env-default:"true"  yaml:"application"`
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
