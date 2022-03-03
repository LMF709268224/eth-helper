package ethevent

import (
	"context"
	"eth-helper/db"
	"eth-helper/erc20"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// 数据库里最小的blockNumber
var blockNumberOfBD = uint64(0)

// 确认数 (可配)
var confirmNum = uint64(10)

// InitTask 初始化交易检查定时器
func InitTask() {
	timeSecond := 60 * 1
	duration := time.Duration(time.Second * time.Duration(timeSecond))

	t := time.NewTicker(duration)
	defer t.Stop()

	for {
		<-t.C

		log.Infoln("task...")
		checkTransfer()
	}
}

func checkTransfer() {
	c := erc20.GetClient()
	defer c.Client.Close()

	// 数据库里最低的高度
	numDB, err := db.GetFristTransfer()
	if err != nil {
		return
	}

	// 链上的高度
	blockNumber, err := getBlockNumber(c)
	if err != nil {
		return
	}

	log.Infof("task...num:%v,bn:%v", numDB, blockNumber)
	if blockNumber-numDB < confirmNum {
		return
	}
	// 如果已经达到确认高度,则开始检查数据库里的交易
	transfers, err := db.GetTransferWithBlocknumber(numDB)
	if err != nil {
		return
	}

	for _, transfer := range transfers {
		if transactionReceipt(c, transfer.Txhash) {
			// TODO 交易完成,发送给mq ,并从表里移除
		} else {
			// TODO 这里要考虑如果某个交易一直没完成,是否要从表里移除,不然会卡住
		}
	}
}

func transactionReceipt(c *erc20.ReturnClient, hash string) bool {
	srt, err := c.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		log.Errorf("transactionReceipt TransactionReceipt err : %v", err)
		return false
	}
	log.Infof("srt------------Status:%v,BlockNumber:%v", srt.Status, srt.BlockNumber)

	return srt.Status == 1
}

// getBlockNumber 获取区块高度
func getBlockNumber(c *erc20.ReturnClient) (uint64, error) {
	bn, err := c.Client.BlockNumber(context.Background())
	if err != nil {
		log.Errorf("getBlockNumber BlockNumber err : %v", err)
	}

	return bn, err
}
