package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Host         string        `env-default:"localhost" yaml:"dbHost"            env:"DB_HOST"`
	Port         int           `env-default:"5432"      yaml:"dbPort"            env:"DB_PORT"`
	User         string        `env-default:"postgres"  yaml:"dbUser"            env:"DB_USER"`
	Password     string        `env-default:"postgres"  yaml:"dbPassword"        env:"DB_PASSWORD"`
	Database     string        `env-default:"postgres"  yaml:"dbDatabase"        env:"DB_DATABASE"`
	SSLMode      string        `env-default:"disable"   yaml:"dbSslmode"         env:"DB_SSLMODE"`
	MaxOpenConns int           `env-default:"25"        yaml:"dbMaxOpenConnects" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int           `env-default:"25"        yaml:"dbMaxIdleConns"    env:"DB_MAX_IDLE_CONNS"`
	MaxTimeLife  time.Duration `env-default:"24h"       yaml:"dbMaxTimeLife"     env:"DB_MAX_TIME_LIFE"`
}

func (d *Database) DBUrl() string {
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%d sslmode=%s",
		d.User, d.Database, d.Password, d.Host, d.Port, d.SSLMode,
	)
}

type Auth struct {
	JwtSecret string        `env-default:"secret" yaml:"jwtSecret" env:"JWT_SECRET"`
	ExpTime   time.Duration `yaml:"expTime"   env:"EXP_TIME"`
}

type Application struct {
	Port int `env-default:"8000" yaml:"port"`
}
type Config struct {
	Env         string `env-default:"local" yaml:"env"          env:"ENV"`
	StoragePath string `                    yaml:"storage_path" env:"storagePath"`
	Database    Database
	Auth        Auth
	Application Application
}

func MustLoad() *Config {
	var cfg Config
	path := fetchConfigPath()
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic("config file not found: " + path)
		}
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic("cannot read config: " + err.Error())
		}
	} else {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(err)
		}

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
