package main

import (
	"eth-helper/message"
	"eth-helper/server"
	"fmt"
	"log"
	"os"

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
	}

	app.Action = func(c *cli.Context) error {
		port := c.String("port")

		message.Init()
		// test
		message.TestHttp()

		// 4、开启Http服务
		params := fmt.Sprintf(":%s", port)
		server.StartHTTPServer(params)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
