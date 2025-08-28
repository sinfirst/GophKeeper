package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
	"github.com/sinfirst/GophKeeper/internal/models"
)

var (
	TokenSetting models.TokenSettings = models.TokenSettings{
		TokenExp:  time.Hour * 12,
		SecretKey: "supersecretkey",
	}
	VersionBuild models.VersionBuild = models.VersionBuild{
		Version: "0.1 beta",
		Date:    "28.08.2025",
	}
)

type Config struct {
	Host        string `env:"BASE_HOST"`
	DatabaseDsn string `env:"DATABASE_URI"`
}

func NewConfig() Config {
	var conf Config
	err := env.Parse(&conf)

	if err != nil {
		fmt.Println(err)
	}

	if conf.DatabaseDsn == "" {
		flag.StringVar(&conf.DatabaseDsn, "d", "postgres://postgres:qwerty12345@localhost:5432/postgres", "database dsn")
	}

	flag.StringVar(&conf.Host, "a", ":3200", "host")

	flag.Parse()

	return conf
}
