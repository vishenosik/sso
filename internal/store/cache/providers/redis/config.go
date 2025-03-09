package redis

import "github.com/blacksmith-vish/sso/internal/lib/config/info"

func init() {
	info.AddToSchema(&Config{})
}

type Config struct {
	Options Options `yaml:"options"`
}

type Options struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func (c *Config) Schema() []byte {
	return append(
		[]byte("Redis:\n"),
		info.ConfigInfoTags(c)...,
	)
}
