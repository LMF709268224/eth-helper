package main

import (
	"eth-helper/config"
	"eth-helper/db"
	"eth-helper/ethevent"
	"eth-helper/server"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var createAddressCommand = &cli.Command{
	Name:    "ca",
	Aliases: []string{"ca"},
	Usage:   "创建地址",
	Action: func(c *cli.Context) error {
		num := c.Args().Get(0)

		n, err := strconv.Atoi(num)
		if err != nil {
			n = 1
		}

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
		// 有参数则用参数，没参数才会使用环境变量
		&cli.StringFlag{
			Name:  "port",
			Value: "7888",
			Usage: "http端口",
			// Destination: &port,
			EnvVars: []string{"BLOCKCHAIN_PORT"},
		},
		&cli.StringFlag{
			Name:  "sqluse",
			Value: "gouse",
			Usage: "mysql用户名",
			// Destination: &sqluse,
			EnvVars: []string{"BLOCKCHAIN_SQLUSE"},
		},
		&cli.StringFlag{
			Name:  "sqlpass",
			Value: "123456",
			Usage: "mysql密码",
			// Destination: &sqlpass,
			EnvVars: []string{"BLOCKCHAIN_SQLPASS"},
		},
		&cli.StringFlag{
			Name:  "sqldab",
			Value: "test",
			Usage: "mysql数据库名",
			// Destination: &sqldab,
			EnvVars: []string{"BLOCKCHAIN_SQLDATABASE"},
		},
		&cli.StringFlag{
			Name:  "c",
			Value: "config.toml",
			Usage: "配置路径",
			// Destination: &sqldab,
			EnvVars: []string{"BLOCKCHAIN_SQLDATABASE"},
		},
	}

	app.Commands = []*cli.Command{
		createAddressCommand,
	}

	app.Action = func(c *cli.Context) error {
		port := c.String("port")
		sqluse := c.String("sqluse")
		sqlpass := c.String("sqlpass")
		sqldatabase := c.String("sqldab")
		configPath := c.String("c")

		// 初始化配置
		config.InitConfig(configPath)

		// 初始化DB
		db.InitDB(sqluse, sqlpass, sqldatabase)

		// 监听交易消息
		ethevent.InitWatchTransfer()

		// 初始化检查交易定时器
		go ethevent.InitTask()

		// 开启Http服务
		params := fmt.Sprintf(":%s", port)
		server.StartHTTPServer(params)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
