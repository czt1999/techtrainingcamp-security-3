package gredis

import (
	"fmt"
	"gin/pkg/settings"
	"time"

	"github.com/go-redis/redis"
)

var rc *redis.Client

// Setup Initialize the Redis instance
func Setup() {

	rc = redis.NewClient(&redis.Options{
		Addr:     settings.RedisSetting.Addr,
		Password: settings.RedisSetting.Password,
		DB:       settings.RedisSetting.DB,
		// Dual
		IdleTimeout: settings.RedisSetting.IdleTimeout,
	})

	if err := rc.Ping().Err(); err != nil {
		panic(err)
	}
}

// Set setting a key/value
func Set(key string, value interface{}, exp time.Duration) error {
	return rc.Set(key, value, exp).Err()
}

func Exist(key string) (bool, error) {
	val, err := rc.Exists(key).Result()
	return val > 0, err
}

func Get(key string) (string, error) {
	val, err := rc.Get(key).Result()
	return val, err
}

func Delete(key string) error {
	return rc.Del(key).Err()
}

func PutWindow(key string, timestampNow int64, windowSize int64) (int, error) {
	pipe := rc.TxPipeline()
	pipe.ZAdd(key, redis.Z{Member: timestampNow, Score: float64(timestampNow)})
	pipe.ZRemRangeByScore(key, "0", fmt.Sprintf("%v", timestampNow-windowSize*1000))
	pipe.Expire(key, time.Duration(windowSize)*time.Second)
	pipe.ZCard(key)
	cmder, err := pipe.Exec()
	if err != nil {
		return 0, err
	}
	// cmder[3]: ZCard(key)
	cnt, err := cmder[3].(*redis.IntCmd).Result()
	if err != nil {
		return 0, err
	}
	return int(cnt), nil
}

func PutWindowWithValue(key string, value interface{}, timestampNow int64, windowSize int64) (int, error) {
	pipe := rc.TxPipeline()
	pipe.ZAdd(key, redis.Z{Member: value, Score: float64(timestampNow)})
	pipe.ZRemRangeByScore(key, "0", fmt.Sprintf("%v", timestampNow-windowSize*1000))
	pipe.Expire(key, time.Duration(windowSize)*time.Second)
	pipe.ZCard(key)
	cmder, err := pipe.Exec()
	if err != nil {
		return 0, err
	}
	// cmder[3]: ZCard(key)
	cnt, err := cmder[3].(*redis.IntCmd).Result()
	if err != nil {
		return 0, err
	}
	return int(cnt), nil
}