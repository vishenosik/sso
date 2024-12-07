package config

import "github.com/blacksmith-vish/sso/internal/lib/config"

type Cache struct {
	Redis Redis `yaml:"redis"`
}
type Redis struct {
	Options Options `yaml:"options"`
}

type Options struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func (ca Cache) getRedisConfig() config.Redis {
	return config.Redis{
		Options: config.Options{
			User:     ca.Redis.Options.User,
			Password: ca.Redis.Options.Password,
			DB:       ca.Redis.Options.DB,
			Host:     ca.Redis.Options.Host,
			Port:     ca.Redis.Options.Port,
		},
	}
}
