package config

import (
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/pkg/env"
)

func (c *ConfigFS) RedisConfig() (redis.Config, error) {
	var conf struct {
		Redis redis.Config `yaml:"redis"`
	}

	if err := parseFile(c, &conf); err != nil {
		return redis.Config{}, err
	}

	env.ReadEnv(&conf)

	return conf.Redis, nil
}
