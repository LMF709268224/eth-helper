package main

import (
	"eth-helper/config"
	"eth-helper/db"
	"eth-helper/erc20"
	"eth-helper/ethevent"
	"eth-helper/redishelper"
	"eth-helper/server"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var testCommand = &cli.Command{
	Name:    "test",
	Aliases: []string{"test"},
	Usage:   "测试",
	Action: func(c *cli.Context) error {
		configPath := c.Args().Get(0)

		// 初始化配置文件
		err := config.InitConfig(configPath)
		if err != nil {
			log.Errorln("InitConfig err :", err)
			return err
		}

		// 初始化DB
		databaseInfo := config.GetDatabaseConfig()
		db.InitDB(databaseInfo.UserName, databaseInfo.UserPassword, databaseInfo.DatabaseName)

		// 初始化以太坊合约相关配置
		ethClientInfo := config.GetEthClientConfig()
		erc20.Init(ethClientInfo.NodeWss, ethClientInfo.ContractAddress)

		// ethevent.TestGetBlock()

		return nil
	},
}

var syncDBAddressToRedis = &cli.Command{
	Name:    "sa",
	Aliases: []string{"sa"},
	Usage:   "同步db里面的地址到redis",
	Action: func(c *cli.Context) error {
		configPath := c.Args().Get(0)

		// 初始化配置文件
		err := config.InitConfig(configPath)
		if err != nil {
			log.Errorln("InitConfig err :", err)
			return err
		}

		// 初始化DB
		databaseInfo := config.GetDatabaseConfig()
		db.InitDB(databaseInfo.UserName, databaseInfo.UserPassword, databaseInfo.DatabaseName)

		redishelper.InitServerAddress(config.GetRedisServer())

		infos := db.GetAllAddressInfo()
		err = redishelper.AddAddress(infos)
		if err != nil {
			log.Errorln("AddAddress err :", err)
			return err
		}

		return nil
	},
}

var createAddressCommand = &cli.Command{
	Name:    "ca",
	Aliases: []string{"ca"},
	Usage:   "创建地址",
	Action: func(c *cli.Context) error {
		num := c.Args().Get(0)
		configPath := c.Args().Get(1)

		// 创建个数
		n, err := strconv.Atoi(num)
		if err != nil {
			n = 1
		}

		// 初始化配置文件
		err = config.InitConfig(configPath)
		if err != nil {
			log.Errorln("InitConfig err :", err)
			return err
		}

		// 初始化DB
		databaseInfo := config.GetDatabaseConfig()
		db.InitDB(databaseInfo.UserName, databaseInfo.UserPassword, databaseInfo.DatabaseName)

		// 初始化以太坊合约相关配置
		ethClientInfo := config.GetEthClientConfig()
		erc20.Init(ethClientInfo.NodeWss, ethClientInfo.ContractAddress)

		err = ethevent.NewAddresss(n)
		if err != nil {
			log.Errorf("main NewAddresss err: %v", err.Error())
		}

		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "eth-helper"
	app.Usage = "a eth helper server"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "c",
			Value: "config.toml",
			Usage: "配置路径",
			// Destination: &sqldab,
			// EnvVars: []string{"BLOCKCHAIN_SQLDATABASE"},
		},
	}

	app.Commands = []*cli.Command{
		createAddressCommand,
		testCommand,
		syncDBAddressToRedis,
	}

	app.Action = func(c *cli.Context) error {
		configPath := c.String("c")

		// 初始化配置文件
		err := config.InitConfig(configPath)
		if err != nil {
			log.Errorln("InitConfig err :", err)
			return err
		}

		// 初始化DB
		databaseInfo := config.GetDatabaseConfig()
		db.InitDB(databaseInfo.UserName, databaseInfo.UserPassword, databaseInfo.DatabaseName)

		// 初始化redis
		redishelper.InitServerAddress(config.GetRedisServer())

		// 初始化以太坊合约相关配置
		ecInfo := config.GetEthClientConfig()
		erc20.Init(ecInfo.NodeWss, ecInfo.ContractAddress)

		// 监听交易消息
		ethevent.InitWatchTransfer()

		// 初始化检查交易定时器
		go ethevent.InitTransferCheckTask(ecInfo.ConfirmBlockmeta)

		// 初始化扫快定时器
		go ethevent.InitScanningBlockTask(config.GetBlockNumber(), ecInfo.ContractAddress, ecInfo.ChainID, ecInfo.NodeHTTPS)

		// 开启Http服务
		params := fmt.Sprintf(":%s", config.GetPort())
		server.StartHTTPServer(params)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
