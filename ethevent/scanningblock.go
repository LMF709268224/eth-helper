package ethevent

import (
	"context"
	"eth-helper/db"
	"eth-helper/erc20"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var blocknumber int

// InitScanningBlockTask 初始化扫快
func InitScanningBlockTask(num int) {
	blocknumber = num

	timeSecond := 60 * 1
	duration := time.Duration(time.Second * time.Duration(timeSecond))

	t := time.NewTicker(duration)
	defer t.Stop()

	for {
		<-t.C

		log.Infoln("task...")
		scanningBlock()
	}
}

func scanningBlock() {
	c := erc20.GetClient()
	defer c.Client.Close()

	// TODO 读数据库看看数据库里的高度 blocknumberDB + 1
	blocknumberDB := 0

	// 对比配置高度 用值大的来开始扫
	if blocknumberDB > blocknumber {
		blocknumber = blocknumberDB
	}

	// 扫快
	block, err := getBlockByNumber(c, int64(blocknumber))
	if err != nil {
		return
	}

	// 处理订单
	err = readTransactions(block)
	if err != nil {
		return
	}

	// TODO 最新num存DB

	// blocknumber ++
	blocknumber++
}

func readTransactions(block *types.Block) error {
	// 在事务中处理
	mdb := db.GetDBConnection()
	err := mdb.Transaction(func(tx *gorm.DB) error {
		for _, transaction := range block.Transactions() {

			if transaction.To() == nil {
				continue
			}

			to := transaction.To().Hex()
			txHash := transaction.Hash().Hex()

			if to != "myAddress" {
				// 不是我们自己的地址
				continue
			}

			// 看看订单是否已经处理 eth_transferdone_tbs
			_, err := db.GetTransferInfo(tx, txHash)
			if err == nil {
				// 在表里找到了 则不处理
				continue
			}

			// 写到待处理表 eth_transfer_tbs
			err = db.SaveNewTransfer(tx, db.EthTransferTb{
				MTo:    to,
				Txhash: txHash,
				Value:  transaction.Value().String(),
			})
			if err != nil {
				log.Errorf("readTransactions SaveNewTransfer err : %v,hash : %s", err, txHash)
				return err
			}
		}

		return nil
	})

	return err
}

func getBlockByNumber(c *erc20.ReturnClient, num int64) (*types.Block, error) {
	block, err := c.Client.BlockByNumber(context.Background(), big.NewInt(num))
	if err != nil {
		log.Errorf("getBlockByNumber BlockByNumber err : %v", err)
	}

	return block, err
}

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

func getBlockByHash(c *erc20.ReturnClient, hash common.Hash) (*types.Block, error) {
	block, err := c.Client.BlockByHash(context.Background(), hash)
	if err != nil {
		log.Errorf("getBlockByHash BlockByHash err : %v", err)
	}

	return block, err
}
