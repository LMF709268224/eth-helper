package ethevent

import (
	"context"
	"eth-helper/erc20"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	block, err := getBlockByNumber(c, int64(num))
	if err != nil {
		return
	}
	log.Infof("getBlockByNumber : %v", block)

	for _, tx := range block.Transactions() {
		// log.Infof("tx:%v", tx)
		// log.Infof("Hash: %s ,Value:%s, To:%s", tx.Hash().Hex(), tx.Value().String(), tx.To().Hex())
		log.Infoln(tx.Hash().Hex())     // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		log.Infoln(tx.Value().String()) // 10000000000000000
		log.Infoln(tx.To())             // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
		log.Infoln("---------------------")
	}
}

func getBlockHeaderByNumber(c *erc20.ReturnClient, num int64) (*types.Header, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), big.NewInt(num))
	if err != nil {
		log.Errorf("getBlockHeaderByNumber HeaderByNumber err : %v", err)
	}

	return header, err
}

func getBlockByNumber(c *erc20.ReturnClient, num int64) (*types.Block, error) {
	block, err := c.Client.BlockByNumber(context.Background(), big.NewInt(num))
	if err != nil {
		log.Errorf("getBlockByNumber BlockByNumber err : %v", err)
	}

	return block, err
}

func getBlockByHash(c *erc20.ReturnClient, hash common.Hash) (*types.Block, error) {
	block, err := c.Client.BlockByHash(context.Background(), hash)
	if err != nil {
		log.Errorf("getBlockByHash BlockByHash err : %v", err)
	}

	return block, err
}
