package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type fsConfig[Type any] struct {
	path string `env:"CONFIG_PATH" default:"./config.yaml"`
}

func NewFSConfig[Type any](
	path string,
) *fsConfig[Type] {
	return &fsConfig[Type]{
		path: path,
	}
}

func (c *fsConfig[Type]) LoadConfig(container *Type) error {

	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return err
	}

	conf := new(Type)

	if err := cleanenv.ReadConfig(c.path, conf); err != nil {
		return err
	}

	return nil
}
