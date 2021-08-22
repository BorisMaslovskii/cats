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
	HmacJWTSecretString string `env:"JWT_SECRET" envDefault:"JWTSampleSecret"`
	CatsDBType          string `env:"CATSDBTYPE" envDefault:"postgres"` // postgres or mongo
	PostgresHost        string `env:"PG_HOST" envDefault:"127.0.0.1"`
	PostgresPort        string `env:"PG_PORT" envDefault:"5432"`
	PostgresUser        string `env:"PG_USER" envDefault:"postgres"`
	PostgresPassword    string `env:"PG_PASS" envDefault:"pgpass"`
	MongoHost           string `env:"MONGO_HOST" envDefault:"127.0.0.1"`
	MongoPort           string `env:"MONGO_PORT" envDefault:"27017"`
	MongoDatabase       string `env:"MONGO_DB" envDefault:"local"`
	MongoCollection     string `env:"MONGO_COLL" envDefault:"cats"`
}

// NewConfig creates an instance of Config
func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
