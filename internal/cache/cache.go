package cache

import (
	"context"
	"fmt"
	"strconv"

	"docker-first/internal/config"

	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

type Cache interface {
	InitData(ctx context.Context) error
	Get(ctx context.Context) (int64, error)
	Increment(ctx context.Context) (int64, error)
}

func InitCache(cfg *config.Config) (Cache, error) {
	rdb := &Redis{
		rds: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("redis:%s", cfg.CachePort),
			Password: "",
			DB:       0,
		}),
	}

	_, err := rdb.rds.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "error in ping redis")
	}

	return rdb, nil
}

type Redis struct {
	rds *redis.Client
}

const counterKey = "counter"

func (r *Redis) InitData(ctx context.Context) error {
	saved, err := r.rds.Set(ctx, counterKey, 0, 0).Result()
	if err != nil {
		return errors.Wrap(err, "save available currencies")
	}

	if saved != "OK" {
		return errors.New("save no info")
	}

	return nil
}

func (r *Redis) Get(ctx context.Context) (int64, error) {
	counterValueStr, err := r.rds.Get(ctx, counterKey).Result()
	if err != nil {
		return 0, errors.Wrap(err, "save available currencies")
	}

	counterValue, err := strconv.ParseInt(counterValueStr, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "unable parse int")
	}

	return counterValue, nil
}

func (r *Redis) Increment(ctx context.Context) (int64, error) {
	counterValue, err := r.rds.Incr(ctx, counterKey).Result()
	if err != nil {
		return 0, errors.Wrap(err, "save available currencies")
	}

	return counterValue, nil
}
