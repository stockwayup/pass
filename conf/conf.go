package conf

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env       string `json:"env" required:"true"`
	DB        *DB
	DebugMode bool `default:"false" json:"debug_mode"`
	Nats      struct {
		Host   string `json:"host"     required:"true"`
		Queues struct {
			Validation string `json:"validation" default:"validation"`
			Generation string `json:"generation" default:"generation"`
		} `json:"queues"`
	} `json:"nats"`
	Password struct {
		Time    uint32 `default:"1"     json:"time"`
		Memory  uint32 `default:"65536" json:"memory"`
		Threads uint8  `default:"4"     json:"threads"`
		KeyLen  uint32 `default:"32"    json:"key_len"`
	} `json:"password"`
}

type DB struct {
	Master     *Master
	Slave      *Slave
	EnabledSSL bool `json:"enabled_ssl"`
}

// Master config.
type Master struct {
	User     string `json:"user"     required:"true"`
	Name     string `json:"name"     required:"true"`
	Host     string `json:"host"     required:"true"`
	Port     string `default:"5432"  json:"port"`
	Password string `json:"password" required:"true"`
	PoolSize int    `default:"50"    json:"pool_size"`
}

// Slave config.
type Slave struct {
	User     string `json:"user"     required:"true"`
	Name     string `json:"name"     required:"true"`
	Host     string `json:"host"     required:"true"`
	Port     string `default:"5432"  json:"port"`
	Password string `json:"password" required:"true"`
	PoolSize int    `default:"50"    json:"pool_size"`
}

func New() *Config {
	c := &Config{}
	path := os.Getenv("CFG_PATH")

	if path == "" {
		path = "./conf/config.json"
	}

	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(c, path); err != nil {
		log.Fatal().Err(err).Msg("config validation errors")
	}

	return c
}
