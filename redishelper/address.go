package redishelper

import (
	"eth-helper/db"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// AddAddress 添加地址
func AddAddress(infos []db.EthAddressTb) error {
	conn := getConn()
	defer conn.Close()

	args := []interface{}{KeyAddress}
	for _, info := range infos {
		args = append(args, info.AddressLower)
	}

	_, err := conn.Do("SADD", args...)
	if err != nil {
		log.Println("AddAddress, SADD err:", err)
	}

	return err
}

// ExistsAddress 判断地址是否存在
func ExistsAddress(address string) bool {
	conn := getConn()
	defer conn.Close()

	address = strings.ToLower(address)

	existe, err := redis.Bool(conn.Do("SISMEMBER", KeyAddress, address))
	if err != nil {
		log.Println("ExistsAddress, SISMEMBER err:", err)
		return false
	}

	return existe
}
