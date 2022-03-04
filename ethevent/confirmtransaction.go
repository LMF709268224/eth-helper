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
// var blockNumberOfBD = uint64(0)

// 确认数 (可配)
var confirmNum = uint64(10)

// 是否在检查确认中,如果是,则不再发起请求
var isChecking = false

// InitTask 初始化交易检查定时器
func InitTask() {
	timeSecond := 60 * 1
	duration := time.Duration(time.Second * time.Duration(timeSecond))

	t := time.NewTicker(duration)
	defer t.Stop()

	for {
		<-t.C

		log.Infoln("task...")
		if !isChecking {
			checkTransfer()
		}
	}
}

func checkTransfer() {
	isChecking = true

	c := erc20.GetClient()
	defer deferCheck(c)

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

	log.Infof("checkTransfer...num:%v,bn:%v", numDB, blockNumber)
	if blockNumber-numDB < confirmNum {
		return
	}
	// 如果已经达到确认高度,则开始检查数据库里的交易
	transfers, err := db.GetTransferWithBlocknumber(numDB)
	if err != nil {
		return
	}

	for _, transfer := range transfers {

		status, err := transactionReceipt(c, transfer.Txhash)
		if err != nil {
			continue
		}

		log.Infof("checkTransfer...交易:%v,状态是:%b", transfer.Txhash, status)
		if status {
			// TODO 交易完成,发送给mq
		}
		// TODO 交易状态不管完成与否,都不再查询,从表里移除,可能要放到'交易完成表''交易失败表'
		db.DeleteTransfer(transfer.ID)
	}
}

func transactionReceipt(c *erc20.ReturnClient, hash string) (bool, error) {
	srt, err := c.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		log.Errorf("transactionReceipt err : %v", err)
		return false, err
	}
	// log.Infof("srt------------Status:%v,BlockNumber:%v", srt.Status, srt.BlockNumber)

	return srt.Status == 1, nil
}

// getBlockNumber 获取区块高度
func getBlockNumber(c *erc20.ReturnClient) (uint64, error) {
	bn, err := c.Client.BlockNumber(context.Background())
	if err != nil {
		log.Errorf("getBlockNumber BlockNumber err : %v", err)
	}

	return bn, err
}

func deferCheck(c *erc20.ReturnClient) {
	c.Client.Close()
	isChecking = false
}
