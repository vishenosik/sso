package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

type Cache struct {
	Redis Redis `yaml:"redis"`
}

type Redis struct {
	Options Options `yaml:"options"`
}

type Options struct {
	User     string
	Password string
	DB       int
	Host     string
	Port     int
}

type Config struct {
	Env       string `validate:"oneof=dev prod test"`
	StorePath string `validate:"required"`
	Cache     Cache  `yaml:"cache"`
}

type ConfigProvider interface {
	FetchConfig() *Config
}

func main() {
	conf := MustLoad()
	fmt.Println(conf)
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
	var res string

	flag.StringVar(&res, "config", "", "path to config file")

	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
