package main

import (
	"eth-helper/db"
	"eth-helper/ethevent"
	"eth-helper/server"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "eth-helper"
	app.Usage = "a eth helper server"

	app.Flags = []cli.Flag{
		// 有参数则用参数，没参数才会使用环境变量
		&cli.StringFlag{
			Name:  "port",
			Value: "7888",
			Usage: "port of the server",
			// Destination: &port,
			EnvVars: []string{"BLOCKCHAIN_PORT"},
		},
		&cli.StringFlag{
			Name:  "sqluse",
			Value: "gouse",
			Usage: "sql username",
			// Destination: &sqluse,
			EnvVars: []string{"BLOCKCHAIN_SQLUSE"},
		},
		&cli.StringFlag{
			Name:  "sqlpass",
			Value: "123456",
			Usage: "sql password",
			// Destination: &sqlpass,
			EnvVars: []string{"BLOCKCHAIN_SQLPASS"},
		},
		&cli.StringFlag{
			Name:  "sqldab",
			Value: "test",
			Usage: "sql database",
			// Destination: &sqldab,
			EnvVars: []string{"BLOCKCHAIN_SQLDATABASE"},
		},
		// TEST
		&cli.StringFlag{
			Name:  "ca",
			Value: "0",
			Usage: "create address num",
			// Destination: &sqldab,
			EnvVars: []string{"BLOCKCHAIN_SQLDATABASE"},
		},
	}

	app.Action = func(c *cli.Context) error {
		port := c.String("port")
		sqluse := c.String("sqluse")
		sqlpass := c.String("sqlpass")
		sqldatabase := c.String("sqldab")

		// 初始化DB
		db.InitDB(sqluse, sqlpass, sqldatabase)

		// test create address
		ca := c.Int("ca")
		if ca > 0 {
			err := ethevent.NewAddresss(ca)
			if err != nil {
				log.Errorf("main NewAddresss err: %v", err.Error())
			}
		}

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
