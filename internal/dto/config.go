package dto

import (
	"fmt"
	"time"
)

type AuthConfig struct {
	SecretAccess  string        `yaml:"SecretAccess"  env:"SECRET_ACCESS"`
	ExpAccess     time.Duration `yaml:"ExpAccess"     env:"EXP_ACCESS"     env-default:"1h"`
	SecretRefresh string        `yaml:"SecretRefresh" env:"SECRET_REFRESH"`
	ExpRefresh    time.Duration `yaml:"ExpRefresh"    env:"EXP_REFRESH"    env-default:"24h"`
}

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

func (d *Database) Url() string {
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%d sslmode=%s",
		d.User, d.Database, d.Password, d.Host, d.Port, d.SSLMode,
	)
}

type Application struct {
	Port int    `env-default:"8000"  yaml:"port"`
	Env  string `env-default:"local" yaml:"env"  env:"ENV"`
}
