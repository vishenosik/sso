package apps

import "github.com/blacksmith-vish/sso/internal/store/models"

type Apps struct {
	cache models.CacheProvider
}

func NewAppsCache(cache models.CacheProvider) *Apps {
	return &Apps{
		cache: cache,
	}
}
