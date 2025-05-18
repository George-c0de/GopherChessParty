package config

import (
	"flag"
	"os"

	"GopherChessParty/internal/dto"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database    dto.Database
	Auth        dto.AuthConfig
	Application dto.Application
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
