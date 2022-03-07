package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type tomlConfig struct {
	Title string
	DB    database `toml:"database"`
}

// InitConfig 初始化配置
func InitConfig(configPath string) {
	var config tomlConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Errorln(err)
		return
	}

	log.Printf("Title: %s\n", config.Title)
	log.Printf("Database: %s %v (Max conn. %d), Enabled? %v\n",
		config.DB.Server, config.DB.Ports, config.DB.ConnMax,
		config.DB.Enabled)
}
