// Package config provides functionality to get the Config struct with all required env variables inside
package config

import (
	"github.com/caarlos0/env/v6"
)

// hmacJWTSecret: for HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent.

// Config provides access for the environment variables
type Config struct {
	HTTPPort            string `env:"HTTP_PORT" envDefault:"1323"`
	HmacJWTSecretString string `env:"JWT_SECRET,notEmpty"`
	PostgresHost        string `env:"PG_HOST" envDefault:"host.docker.internal"`
	PostgresPort        string `env:"PG_PORT" envDefault:"5432"`
	PostgresUser        string `env:"PG_USER,required"`
	PostgresPassword    string `env:"PG_PASS,required"`
	MongoURI            string `env:"MONGO_URI"`
	MongoDatabase       string `env:"MONGO_DB"`
	MongoCollection     string `env:"MONGO_COLL"`
}

// NewConfig creates an instance of Config
func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
