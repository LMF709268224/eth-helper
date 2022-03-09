package redishelper

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	// KeyBlockNumber block number
	KeyBlockNumber = "eth:bn"
	// KeyAddress block number
	KeyAddress = "eth:address"
)

var (
	pool *redis.Pool

	once sync.Once
)

var address string

// newPool 新建redis连接池
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", address) },
	}
}

// InitServerAddress init
func InitServerAddress(addr string) {
	address = addr
}

// getConn get a redis connection
func getConn() redis.Conn {
	if pool == nil {
		once.Do(func() {
			pool = newPool()
		})
	}

	return pool.Get()
}
