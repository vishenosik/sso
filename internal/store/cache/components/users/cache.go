package users

import "github.com/blacksmith-vish/sso/internal/store/models"

type Users struct {
	cache models.CacheProvider
}

func NewUsersCache(cache models.CacheProvider) *Users {
	return &Users{
		cache: cache,
	}
}
