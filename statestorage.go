package discordgo

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

// StateCache represents an arbitrary state cache
type StateCache interface {
	Get(set, key string) (interface{}, error)
	Put(set, key string, value interface{}) error
	Del(set, key string) (interface{}, error)
	DelSet(set string) error
}

// RedisAdapter is an implementation of Cache for redis
type RedisAdapter struct {
	client *redis.Client
}

// NewRedisAdapter returns a new instance of RedisAdapter
func NewRedisAdapter(redisOptions interface{}) RedisAdapter {
	options := redisOptions.(redis.Options)

	return RedisAdapter{client: redis.NewClient(&options)}
}

// Get returns the cache value at the given key
func (c RedisAdapter) Get(set, key string) (interface{}, error) {
	res, err := c.client.HGet(set, key).Result()
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal([]byte(res), value)

	return value, err
}

// Put inserts a value into the cache at the given key
func (c RedisAdapter) Put(set, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = c.client.HSet(set, key, data).Result()

	return err
}

// Del removes the given key and returns the value
func (c RedisAdapter) Del(set, key string) (interface{}, error) {
	value, err := c.Get(set, key)
	if err != nil {
		return nil, err
	}

	_, err = c.client.HDel(set, key).Result()

	return value, err
}

// DelSet removes an entire sub cache from the state cache
func (c RedisAdapter) DelSet(set string) error {
	_, err := c.client.Del(set).Result()
	return err
}

// LocalAdapter is an implementation of Cache for local memory
type LocalAdapter struct {
	cache *cache.Cache
}

// NewLocalAdapter returns a new instance of LocalAdapter
func NewLocalAdapter() LocalAdapter {
	return LocalAdapter{cache: cache.New(cache.NoExpiration, 10*time.Minute)}
}

// Get returns the cache value at the given key
func (c LocalAdapter) Get(set, key string) (interface{}, error) {
	cacheSet, ok := c.cache.Get(set)
	if !ok {
		return nil, errors.New("cache set not found")
	}

	setValue := cacheSet.(cache.Cache)

	value, ok := setValue.Get(key)
	if !ok {
		return nil, errors.New("value does not exist")
	}

	return value, nil
}

// Put inserts a value into the cache at the given key
func (c LocalAdapter) Put(set, key string, value interface{}) error {
	var cacheSet *cache.Cache
	s, ok := c.cache.Get(set)
	if !ok {
		cacheSet = cache.New(cache.NoExpiration, 10*time.Minute)
		c.cache.Set(set, cacheSet, 0)
	} else {
		cacheSet = s.(*cache.Cache)
	}

	cacheSet.Set(key, value, 0)

	return nil
}

// Del removes the given key and returns the value
func (c LocalAdapter) Del(set, key string) (interface{}, error) {
	var cacheSet *cache.Cache
	s, ok := c.cache.Get(set)
	if !ok {
		return nil, errors.New("cache set does not exist")
	}

	cacheSet = s.(*cache.Cache)
	value, ok := cacheSet.Get(key)
	if !ok {
		return nil, errors.New("value not in cache")
	}

	cacheSet.Delete(key)
	return value, nil
}

// DelSet removes an entire sub cache from the state cache
func (c LocalAdapter) DelSet(set string) error {
	c.cache.Delete(set)
	return nil
}
