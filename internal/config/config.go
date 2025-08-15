package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

var TokenExp = time.Hour * 12
var SecretKey = "supersecretkey"

type Config struct {
	Host        string `env:"BASE_URL"`
	DatabaseDsn string `env:"DATABASE_URI"`
}

func NewConfig() Config {
	var conf Config
	err := env.Parse(&conf)

	if err != nil {
		fmt.Println(err)
	}

	if conf.DatabaseDsn == "" {
		flag.StringVar(&conf.DatabaseDsn, "d", "postgres://postgres:1@localhost:5432/postgres", "database dsn")
	}

	flag.StringVar(&conf.Host, "a", ":3200", "host")

	flag.Parse()

	return conf
}
