package cache

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// CatalogCache caches catalog JSON blobs (products list responses).
type CatalogCache interface {
	GetJSON(key string, dest any) bool
	SetJSON(key string, v any, ttl time.Duration)
	InvalidateCatalog()
}

type noopCache struct{}

func (noopCache) GetJSON(string, any) bool            { return false }
func (noopCache) SetJSON(string, any, time.Duration) {}
func (noopCache) InvalidateCatalog()                 {}

type memoryCache struct {
	mu   sync.Mutex
	data map[string]memEntry
}

type memEntry struct {
	raw []byte
	exp time.Time
}

func NewMemoryCatalogCache() CatalogCache {
	return &memoryCache{data: make(map[string]memEntry)}
}

func (m *memoryCache) GetJSON(key string, dest any) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.data[key]
	if !ok || time.Now().After(e.exp) {
		return false
	}
	return json.Unmarshal(e.raw, dest) == nil
}

func (m *memoryCache) SetJSON(key string, v any, ttl time.Duration) {
	raw, err := json.Marshal(v)
	if err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = memEntry{raw: raw, exp: time.Now().Add(ttl)}
}

func (m *memoryCache) InvalidateCatalog() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k := range m.data {
		if strings.HasPrefix(k, "catalog:") {
			delete(m.data, k)
		}
	}
}

type redisCache struct {
	rdb *redis.Client
}

func NewRedis(addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

func NewRedisCatalogCache(rdb *redis.Client) CatalogCache {
	return &redisCache{rdb: rdb}
}

func (r *redisCache) GetJSON(key string, dest any) bool {
	s, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil || s == "" {
		return false
	}
	return json.Unmarshal([]byte(s), dest) == nil
}

func (r *redisCache) SetJSON(key string, v any, ttl time.Duration) {
	raw, err := json.Marshal(v)
	if err != nil {
		return
	}
	_ = r.rdb.Set(context.Background(), key, raw, ttl).Err()
	_ = r.rdb.SAdd(context.Background(), "catalog:index", key).Err()
}

func (r *redisCache) InvalidateCatalog() {
	ctx := context.Background()
	keys, err := r.rdb.SMembers(ctx, "catalog:index").Result()
	if err != nil {
		return
	}
	if len(keys) > 0 {
		_ = r.rdb.Del(ctx, keys...).Err()
	}
	_ = r.rdb.Del(ctx, "catalog:index").Err()
}

// NoOp returns a cache that never hits.
func NoOp() CatalogCache { return noopCache{} }
