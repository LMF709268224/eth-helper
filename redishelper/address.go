package redishelper

import (
	"eth-helper/db"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// AddAddress
func AddAddress(infos []db.EthAddressTb) error {
	conn := getConn()
	defer conn.Close()

	args := []interface{}{KeyAddress}
	for _, info := range infos {
		lower := strings.ToLower(info.Address)
		args = append(args, lower)
	}

	_, err := conn.Do("SADD", args...)
	if err != nil {
		log.Println("AddAddress, SADD err:", err)
	}

	return err
}

// ExistsAddress
func ExistsAddress(address string) bool {
	conn := getConn()
	defer conn.Close()

	existe, err := redis.Bool(conn.Do("SISMEMBER", KeyAddress, address))
	if err != nil {
		log.Println("ExistsAddress, SISMEMBER err:", err)
		return false
	}

	return existe
}
