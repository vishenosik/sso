package config

import (
	"flag"
	"os"

	"log"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/joho/godotenv"

	"github.com/ilyakaznacheev/cleanenv"
)

func init() {

	flag.StringVar(&res, "config", "", "path to config file")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var res string

type Config struct {
	Env       string   `yaml:"env"`
	StorePath string   `yaml:"store_path"`
	Services  Services `yaml:"services"`
	Servers   Servers  `yaml:"servers"`
	Cache     Cache    `yaml:"cache"`
}

func (conf Config) FetchConfig() *config.Config {
	return &config.Config{
		Env:                   conf.env(),
		StorePath:             conf.storePath(),
		AuthenticationService: conf.Services.getAuthenticationService(),
		GrpcConfig:            conf.Servers.getGrpcConfig(),
		RestConfig:            conf.Servers.getRestConfig(),
		Redis:                 conf.Cache.getRedisConfig(),
	}
}

func MustLoad() *Config {

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	conf := new(Config)

	if err := cleanenv.ReadConfig(path, conf); err != nil {
		panic("failed to parse config file: " + err.Error())
	}

	return conf
}

// flag > env > default
func fetchConfigPath() string {
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (conf Config) storePath() string {
	if conf.StorePath == "" {
		return os.Getenv("STORE_PATH")
	}
	return conf.StorePath
}

func (conf Config) env() string {
	if conf.Env == "" {
		return os.Getenv("ENV")
	}
	return conf.Env
}
