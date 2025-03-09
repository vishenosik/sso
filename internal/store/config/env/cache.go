package env

import (
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/pkg/env"
)

func (c *envConfig) RedisConfig() (redis.Config, error) {
	var Redis redis.Config
	env.ReadEnv(&Redis)
	return Redis, nil
}
