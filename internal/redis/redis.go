package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"carbonite/client/internal/data"

	redis "github.com/redis/go-redis/v9"
	"github.com/rs/xid"
)

const hash = "carbonit:sessions:"

type (
	// Client is the Redis client structure
	Client struct {
		*redis.Client
	}
)

func getHashID(key string) string {
	return hash + key
}

func (c Client) RemoveToken(ID string) error {
	_, err := c.LoadToken(ID)
	if err != nil {
		return err
	}

	ctx := context.Background()
	c.Del(ctx, getHashID(ID))
	return nil
}

func (c Client) LoadToken(ID string) (*data.TokenPayload, error) {
	ctx := context.Background()
	result, err := c.HGet(ctx, getHashID(ID), "token").Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if result == "" {
		return nil, fmt.Errorf("find: not found")
	}

	token := &data.TokenPayload{}
	if err := json.Unmarshal([]byte(result), token); err != nil {
		return nil, fmt.Errorf("find: unmarshal error: %w", err)
	}

	return token, nil
}

func (c Client) ReplaceToken(ID string, token *data.TokenPayload) error {
	c.RemoveToken(ID)

	ctx := context.Background()
	data, _ := json.Marshal(token)
	if _, err := c.HSetNX(ctx, getHashID(ID), "token", data).Result(); err != nil {
		return fmt.Errorf("create: redis error: %w", err)
	}
	c.Expire(ctx, getHashID(ID), time.Hour*24)

	return nil
}

func (c Client) SaveToken(token *data.TokenPayload) (string, error) {
	ID := xid.New().String()
	ctx := context.Background()
	data, _ := json.Marshal(token)
	if _, err := c.HSetNX(ctx, getHashID(ID), "token", data).Result(); err != nil {
		return "", fmt.Errorf("create: redis error: %w", err)
	}
	c.Expire(ctx, getHashID(ID), time.Hour*24)

	return ID, nil
}

func New(URL string) (*Client, error) {
	options, err := redis.ParseURL(URL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
		/*	TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},*/
	})

	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
