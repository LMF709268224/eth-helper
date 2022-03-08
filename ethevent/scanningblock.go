package ethevent

import (
	"context"
	"eth-helper/erc20"
	"math/big"

	log "github.com/sirupsen/logrus"
)

// TestGetBlock 测试
func TestGetBlock() {
	c := erc20.GetClient()
	defer c.Client.Close()

	num, err := getBlockNumber(c)
	if err != nil {
		return
	}
	log.Infoln("getBlockNumber num :", num)

	getBlockHeaderByNumber(c, int64(num))
}

func getBlockHeaderByNumber(c *erc20.ReturnClient, num int64) {
	header, err := c.Client.HeaderByNumber(context.Background(), big.NewInt(num))
	if err != nil {
		log.Errorf("getBlockHeaderByNumber HeaderByNumber err : %v", err)
	} else {
		log.Infof("header : %v ", header)
	}
}
