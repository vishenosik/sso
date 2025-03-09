package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var configPath string

func init() {

	flag.StringVar(&configPath, "config.file", "", "path to config file")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Decoder interface {
	Decode(v any) error
}

type OsFile interface {
	io.Reader
	Close() error
}

type ConfigFS struct {
	decoder Decoder
	path    string
	fs      OsFile
}

func MustLoad() *ConfigFS {

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		panic(err)
	}

	cfs := &ConfigFS{
		path: path,
		fs:   file,
	}

	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".yaml", ".yml":
		cfs.decoder = yaml.NewDecoder(file)
	case ".json":
		cfs.decoder = json.NewDecoder(file)
	default:
		panic(fmt.Errorf("file format '%s' doesn't supported by the parser", ext))
	}

	return cfs
}

// flag > env > default
func fetchConfigPath() string {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}

func (cfs *ConfigFS) Close() error {
	return cfs.fs.Close()
}

func parseFile[Type any](cfs *ConfigFS, container *Type) error {
	if err := cfs.decoder.Decode(container); err != nil {
		return fmt.Errorf("config file parsing error: %s", err.Error())
	}
	return nil
}
