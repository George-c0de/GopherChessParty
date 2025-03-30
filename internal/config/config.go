package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Database struct {
	Host         string        `yaml:"DB_HOST"              env-default:"localhost"`
	Port         int           `yaml:"DB_PORT"              env-default:"5432"`
	User         string        `yaml:"DB_USER"              env-default:"postgres"`
	Password     string        `yaml:"DB_PASSWORD"          env-default:"postgres"`
	Database     string        `yaml:"DB_DATABASE"          env-default:"postgres"`
	SSLMode      string        `yaml:"DB_SSLMODE"           env-default:"disable"`
	MaxOpenConns int           `yaml:"DB_MAX_OPEN_CONNECTS" env-default:"25"`
	MaxIdleConns int           `yaml:"DB_MAX_IDLE_CONNS"    env-default:"25"`
	MaxTimeLife  time.Duration `yaml:"DB_MAX_TIME_LIFE"     env-default:"24h"`
}

type Auth struct {
	JwtSecret string        `yaml:"JWT_SECRET" env-default:"secret"`
	ExpTime   time.Duration `yaml:"EXP_TIME"   env-default:"24h"`
}

func (a *Auth) MustParseECDSAPrivateKey() *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(a.JwtSecret))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		panic("failed to parse PEM block containing the EC private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}

func (d *Database) DBUrl() string {
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s  port=%d sslmode=%s",
		d.User, d.Database, d.Password, d.Host, d.Port, d.SSLMode,
	)
}

type Config struct {
	Env         string   `yaml:"env"          env-default:"local"`
	StoragePath string   `yaml:"storage_path"                     env-required:"true"`
	Database    Database `yaml:"database"                         env-required:"true"`
	Auth        Auth     `yaml:"auth"                             env-required:"true"`
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
