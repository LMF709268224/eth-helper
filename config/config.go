package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

var conf *TomlConfig

// Database 数据库配置
type Database struct {
	UserName     string `toml:"user_name"`
	UserPassword string `toml:"user_password"`
	DatabaseName string `toml:"database_name"`
}

// Client Client链配置
type Client struct {
	ConfirmBlockmeta uint64 `toml:"confirm_blockmeta"`
	NodeWss          string `toml:"node_wss"`
	ContractAddress  string `toml:"contract_address"`
	NodeHTTPS        string `toml:"node_https"`
	ChainID          int    `toml:"chainid"`
}

// TomlConfig 配置
type TomlConfig struct {
	Port        string   `toml:"http_port"`
	BlockNumber int      `toml:"block_number"`
	RedisServer string   `toml:"redis_server"`
	DB          Database `toml:"database"`
	EC          Client   `toml:"eth_client"`
	PC          Client   `toml:"polygon_client"`
}

// InitConfig 初始化配置
func InitConfig(configPath string) error {
	var config TomlConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Errorln(err)
		return err
	}

	conf = &config

	return nil
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() Database {
	return conf.DB
}

// GetEthClientConfig 获取以太坊链配置
func GetEthClientConfig() Client {
	return conf.EC
}

// GetPolygonClientConfig 获取Polygon链配置
func GetPolygonClientConfig() Client {
	return conf.PC
}

// GetPort 获取http端口
func GetPort() string {
	return conf.Port
}

// GetBlockNumber 获取扫快高度
func GetBlockNumber() int {
	return conf.BlockNumber
}

func GetRedisServer() string {
	return conf.RedisServer
}
