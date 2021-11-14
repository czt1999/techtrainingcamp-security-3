package main

//设置连接池限制连接类

import (
	"time"

	"github.com/go-redis/redis"
)

// 根据redisURL创建新的limiter并返回
func NewLimiter(redisURL string) (*Limiter, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	rc := redis.NewClient(opts)
	if err := rc.Ping().Err(); err != nil {
		return nil, err
	}

	return &Limiter{rc: rc}, nil
}

// 通过redis的value判断第几次访问并返回是否允许访问
func (l *Limiter) Allow(key string, events int64, per time.Duration) bool {
	curr := l.rc.LLen(key).Val()
	if curr >= events {
		return false
	}

	if v := l.rc.Exists(key).Val(); v == 0 {
		pipe := l.rc.TxPipeline()
		pipe.RPush(key, key)
		//设置过期时间
		pipe.Expire(key, per)
		_, _ = pipe.Exec()
	} else {
		l.rc.RPushX(key, key)
	}

	return true
}
