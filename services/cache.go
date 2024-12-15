package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"github.com/yourusername/yourproject/models"
)

var (
	CacheClient *redis.Client
	ctx         = context.Background()
)

func InitCache(redisHost, redisPort string) {
	CacheClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
}

func GetProductByID(id int) (*models.Product, error) {
	val, err := CacheClient.Get(ctx, fmt.Sprintf("product:%d", id)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("product not found in cache")
	} else if err != nil {
		return nil, err
	}

	var product models.Product
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func SetProductByID(id int, product models.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = CacheClient.Set(ctx, fmt.Sprintf("product:%d", id), data, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func InvalidateProductCache(id int) error {
	err := CacheClient.Del(ctx, fmt.Sprintf("product:%d", id)).Err()
	if err != nil {
		return err
	}

	return nil
}
