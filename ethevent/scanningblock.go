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

	// header, err := getBlockHeaderByNumber(c, int64(num))
	// if err != nil {
	// 	return
	// }

	// block, err := getBlockByHash(c, header.TxHash)
	// if err != nil {
	// 	return
	// }
	// log.Infof("getBlockByHash : %v", block)

	block, err := getBlockByNumber(c, int64(num))
	if err != nil {
		return
	}
	log.Infof("getBlockByNumber : %v", block)

	for _, tx := range block.Transactions() {
		log.Infof("Hash: %s ,Value:%s, To:%s", tx.Hash().Hex(), tx.Value().String(), tx.To().Hex())

		// chainID, err := c.Client.NetworkID(context.Background())
		// if err != nil {
		// 	log.Errorf("NetworkID:%v", err)
		// 	continue
		// }

		// if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
		// 	fmt.Println(msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
		// }
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
