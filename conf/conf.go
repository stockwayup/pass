package conf

import (
	"os"

	pubsub "github.com/soulgarden/rmq-pubsub"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env       string `json:"env" required:"true"`
	DB        *DB
	DebugMode bool `json:"debug_mode"  default:"false"`
	RMQ       struct {
		Host     string `json:"host"      required:"true"`
		Port     string `json:"port"      default:"6379"`
		User     string `json:"user"      default:"guest"`
		Password string `json:"password"  default:"guest"`

		Queues struct {
			GenerateIn  *pubsub.Cfg `json:"generate_in"`
			GenerateOut *pubsub.Cfg `json:"generate_out"`
			ValidateIn  *pubsub.Cfg `json:"validate_in"`
			ValidateOut *pubsub.Cfg `json:"validate_out"`
		} `json:"queues"`
	} `json:"rmq"`
	Password struct {
		Time    uint32 `json:"time"    default:"1"`
		Memory  uint32 `json:"memory"  default:"65536"`
		Threads uint8  `json:"threads" default:"4"`
		KeyLen  uint32 `json:"key_len" default:"32"`
	} `json:"password"`
}

type DB struct {
	Master     *Master
	Slave      *Slave
	EnabledSSL bool `json:"enabled_ssl"`
}

// Master config.
type Master struct {
	User     string `json:"user"      required:"true"`
	Name     string `json:"name"      required:"true"`
	Host     string `json:"host"      required:"true"`
	Port     string `json:"port"      default:"5432"`
	Password string `json:"password"  required:"true"`
	PoolSize int    `json:"pool_size" default:"50"`
}

// Slave config.
type Slave struct {
	User     string `json:"user"      required:"true"`
	Name     string `json:"name"      required:"true"`
	Host     string `json:"host"      required:"true"`
	Port     string `json:"port"      default:"5432"`
	Password string `json:"password"  required:"true"`
	PoolSize int    `json:"pool_size" default:"50"`
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
