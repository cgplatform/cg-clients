package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type MongoConfig struct {
	Host     string `env:"MONGO_HOST"`
	User     string `env:"MONGO_USER"`
	Password string `env:"MONGO_PASSWORD"`
	Database string `env:"MONGO_DATABASE"`
}

type HTTPConfig struct {
	Address string `env:"HTTP_ADDRESS"`
	Port    string `env:"HTTP_PORT"`
}

type MailConfig struct {
	Domain string `env:"MAIL_DOMAIN"`
	Key    string `env:"MAIL_KEY"`
}

var (
	Mongo = MongoConfig{}
	HTTP  = HTTPConfig{}
	Mail  = MailConfig{}
)

func Load() error {
	err := godotenv.Load()

	if err != nil {
		return err
	}

	opts := env.Options{RequiredIfNoDef: true}

	if err := env.Parse(&Mongo, opts); err != nil {
		return err
	}

	if err := env.Parse(&HTTP, opts); err != nil {
		return err
	}

	if err := env.Parse(&Mail, opts); err != nil {
		return err
	}

	return nil
}
