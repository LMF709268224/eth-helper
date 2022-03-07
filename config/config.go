package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

var conf *TomlConfig

type Database struct {
	UserName     string `toml:"user_name"`
	UserPassword string `toml:"user_password"`
	DatabaseName string `toml:"database_name"`
}

type EthClient struct {
	NodeWss         string `toml:"node_wss"`
	ContractAddress string `toml:"contract_address"`
}

type TomlConfig struct {
	Port string    `toml:"http_port"`
	DB   Database  `toml:"database"`
	EC   EthClient `toml:"eth_client"`
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

	// log.Printf("Title: %s\n", config.Port)
	// log.Printf("Database: %s %v (Max conn. %d), \n",
	// 	config.DB.DatabaseName, config.DB.UserName, config.DB.UserPassword)
	// log.Printf("EthClient: %s %v \n",
	// 	config.EC.ContractAddress, config.EC.NodeWss)
}

func GetDatabaseConfig() Database {
	return conf.DB
}

func GetEthClientConfig() EthClient {
	return conf.EC
}

func GetPort() string {
	return conf.Port
}
