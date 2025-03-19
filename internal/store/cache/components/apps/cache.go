package apps

import "github.com/vishenosik/sso/internal/store/models"

type Apps struct {
	cache models.CacheProvider
}

func NewAppsCache(cache models.CacheProvider) *Apps {
	return &Apps{
		cache: cache,
	}
}
