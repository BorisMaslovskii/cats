package config

import (
	"github.com/caarlos0/env/v6"
)

// hmacJWTSecret: for HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent.

type Config struct {
	HmacJWTSecretString string `env:"JWTSECRET,notEmpty"`
	PostgresUser        string `env:"PGUSER,required"`
	PostgresPassword    string `env:"PGPASS,required"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
