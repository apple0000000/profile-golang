package redis

import (
	"context"
	"time"

	"profile-golang/common/models"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient(addr, password string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Client{client: rdb}, nil
}

func (c *Client) Get(key string) (string, error) {
	ctx := context.Background()
	return c.client.Get(ctx, key).Result()
}

func (c *Client) MGet(keys []string) ([]interface{}, error) {
	ctx := context.Background()
	return c.client.MGet(ctx, keys...).Result()
}

func (c *Client) Set(key, value string) error {
	ctx := context.Background()
	return c.client.Set(ctx, key, value, 0).Err()
}

func (c *Client) MSet(items []models.Item) error {
	ctx := context.Background()

	pairs := make([]interface{}, len(items)*2)
	for i, item := range items {
		pairs[i*2] = item.Key
		pairs[i*2+1] = item.Value
	}

	return c.client.MSet(ctx, pairs...).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
