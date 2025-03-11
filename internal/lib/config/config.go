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

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

type Config struct {
	Env                   string `env:"ENV" default:"dev" validate:"oneof=dev prod test" desc:"The environment in which the application is running"`
	StorePath             string `env:"STORE_PATH" default:"./storage/sso.db" validate:"required" desc:"Path to sqlite store"`
	AuthenticationService AuthenticationService
	GrpcConfig            GrpcServer
	RestConfig            RestServer
	Redis                 Redis
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
	Password string `env:"REDIS_PASSWORD" desc:"Redis user's password"`
	DB       int    `env:"REDIS_DB" default:"0" desc:"Redis database connection"`
	Host     string `env:"REDIS_HOST" default:"127.0.0.1" desc:"Redis server host"`
	Port     uint16 `env:"REDIS_PORT" default:"6380" desc:"Redis server port"`
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
