package redishelper

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

// GetBlockNumber 获取已处理高度
func GetBlockNumber(chainID int) int {
	conn := getConn()
	defer conn.Close()

	num, err := redis.Int(conn.Do("HGET", KeyBlockNumber, chainID))
	if err != nil {
		log.Println("GetBlockNumber, HGET err:", err)
		return 0
	}

	return num
}

// SaveBlockNumber 保存已处理高度
func SaveBlockNumber(chainID int, num int) error {
	conn := getConn()
	defer conn.Close()

	_, err := conn.Do("HMSET", KeyBlockNumber, chainID, num)
	if err != nil {
		log.Println("SaveBlockNumber, HMSET err:", err)
	}

	return err
}
