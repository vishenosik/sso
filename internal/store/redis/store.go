package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/blacksmith-vish/sso/internal/domain"
	cache "github.com/blacksmith-vish/sso/internal/lib/cache/redis"
	"github.com/blacksmith-vish/sso/internal/store/models"
)

type CachedAuthentication struct {
	store models.AuthenticationStore
	cache *cache.RedisCache
}

func NewCachedAuthentication(store models.AuthenticationStore, cache *cache.RedisCache) *CachedAuthentication {
	return &CachedAuthentication{
		store: store,
		cache: cache,
	}
}

func (ca *CachedAuthentication) Register(ctx context.Context, email, password string) error {
	err := ca.store.Register(ctx, email, password)
	if err != nil {
		return err
	}

	// Clear any cached user data for this email
	ca.cache.Delete(ctx, fmt.Sprintf("user:%s", email))

	return nil
}

func (ca *CachedAuthentication) Login(ctx context.Context, email, password string) (string, error) {
	token, err := ca.store.Login(ctx, email, password)
	if err != nil {
		return "", err
	}

	// Cache the token
	err = ca.cache.Set(ctx, fmt.Sprintf("token:%s", email), token, 24*time.Hour)
	if err != nil {
		// Log the error, but don't fail the login
		fmt.Printf("Failed to cache token: %v\n", err)
	}

	return token, nil
}

func (ca *CachedAuthentication) Logout(ctx context.Context, token string) error {
	err := ca.store.Logout(ctx, token)
	if err != nil {
		return err
	}

	// Clear the cached token
	ca.cache.Delete(ctx, fmt.Sprintf("token:%s", token))

	return nil
}

func (ca *CachedAuthentication) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user models.User
	cacheKey := fmt.Sprintf("user:%s", email)

	// Try to get user from cache
	err := ca.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// If not in cache, get from store
	user, err = ca.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Cache the user
	err = ca.cache.Set(ctx, cacheKey, user, 1*time.Hour)
	if err != nil {
		// Log the error, but don't fail the operation
		fmt.Printf("Failed to cache user: %v\n", err)
	}

	return &user, nil
}

func (ca *CachedAuthentication) ValidateToken(ctx context.Context, token string) (bool, error) {
	// Try to get from cache first
	var isValid bool
	cacheKey := fmt.Sprintf("valid_token:%s", token)

	err := ca.cache.Get(ctx, cacheKey, &isValid)
	if err == nil {
		return isValid, nil
	}

	// If not in cache, validate from store
	isValid, err = ca.store.ValidateToken(ctx, token)
	if err != nil {
		return false, err
	}

	// Cache the result
	err = ca.cache.Set(ctx, cacheKey, isValid, 15*time.Minute)
	if err != nil {
		// Log the error, but don't fail the operation
		fmt.Printf("Failed to cache token validation: %v\n", err)
	}

	return isValid, nil
}
