package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/blacksmith-vish/sso/pkg/collections"
	"github.com/blacksmith-vish/sso/pkg/env"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var (
	conf *Config
	//
	ErrServerPortMustBeUnique = errors.New("port numbers must be unique")
)

type Config struct {
	Env                   string `env:"ENV" default:"dev" validate:"oneof=dev prod test" desc:"The environment in which the application is running"`
	StorePath             string `env:"STORE_PATH" default:"./storage/sso.db" validate:"required" desc:"Path to sqlite store"`
	AuthenticationService AuthenticationService
	GrpcConfig            GrpcServer
	RestConfig            RestServer
	Redis                 Redis
	Dgraph                Dgraph
	service
}

type RestServer struct {
	Port uint16 `env:"REST_PORT" default:"8080" desc:"REST server port"`
}

type GrpcServer struct {
	Port uint16 `env:"GRPC_PORT" default:"44844" desc:"gRPC server port"`
}

type AuthenticationService struct {
	TokenTTL time.Duration `env:"AUTHENTICATION_TOKEN_TTL" default:"1h" desc:"Authentication service standart TTL"`
}

type Redis struct {
	User     string `env:"REDIS_USER" desc:"Redis user"`
	Password string `env:"REDIS_USER_PASSWORD" desc:"Redis user's password"`
	DB       int    `env:"REDIS_DB" default:"0" desc:"Redis database connection"`
	Host     string `env:"REDIS_HOST" default:"localhost" desc:"Redis server host"`
	Port     uint16 `env:"REDIS_PORT" default:"6380" desc:"Redis server port"`
}

type Dgraph struct {
	User     string `env:"DGRAPH_USER" desc:"Dgraph user"`
	Password string `env:"DGRAPH_USER_PASSWORD" desc:"Dgraph user's password"`
	GrpcHost string `env:"DGRAPH_GRPC_HOST" default:"localhost" desc:"Dgraph server host"`
	GrpcPort uint16 `env:"DGRAPH_GRPC_PORT" default:"9080" desc:"Dgraph server port"`
}

type Vault struct {
	Address string `env:"VAULT_ADDRESS" default:"localhost:8200" desc:"Vault address"`
	Token   string `env:"VAULT_TOKEN" desc:"Vault token"`
}

type service struct {
	vault Vault
}

func init() {

	flag.BoolFunc("config.info", "Show config schema information", env.ConfigInfo[Config](os.Stdout))
	flag.Func("config.doc", "Update config example in docs", env.ConfigDoc[Config]())

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

func EnvConfig() *Config {
	if conf == nil {
		cfg := env.ReadEnv[Config]()
		conf = &cfg
		if collections.HasDuplicates([]uint16{conf.GrpcConfig.Port, conf.RestConfig.Port}) {
			panic(ErrServerPortMustBeUnique)
		}
	}
	return conf
}
