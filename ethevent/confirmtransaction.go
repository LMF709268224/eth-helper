package ethevent

import (
	"context"
	"eth-helper/db"
	"eth-helper/erc20"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 确认数
var confirmBlockmeta = uint64(2)

// 是否在检查确认中,如果是,则不再发起请求
var isChecking = false

// InitTask 初始化交易检查定时器
func InitTask(cb uint64) {
	confirmBlockmeta = cb

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
	c := erc20.GetClient()
	defer func() {
		c.Client.Close()
		isChecking = false
	}()

	isChecking = true

	// 数据库里最低的高度
	numDB, err := db.GetMinBlocknumber()
	if err != nil {
		return
	}
	// log.Infoln("checkTransfer numDB...", numDB)

	// 链上的高度
	blockNumber, err := getBlockNumber(c)
	if err != nil {
		return
	}

	log.Infof("checkTransfer...num:%v,bn:%v", numDB, blockNumber)

	if blockNumber-numDB < confirmBlockmeta {
		return
	}

	// 如果已经达到确认高度,则开始检查数据库里的交易
	transfers, err := db.GetTransferWithBlocknumber(numDB)
	if err != nil {
		return
	}

	for _, transfer := range transfers {
		// 检查交易状态
		status, err := transactionReceipt(c, transfer.Txhash)

		errMsg := ""
		if err != nil {
			// 如果是 not found
			errMsg = err.Error()
		}

		mdb := db.GetDBConnection()
		err = mdb.Transaction(func(tx *gorm.DB) error {
			if status == 1 {
				// TODO 交易完成,发送给mq
				log.Infof("checkTransfer...ID:%d,交易:%v,完成了,要发送给前端.......", transfer.ID, transfer.Txhash)
			}

			// 记录到交易完成表
			err = db.SaveTransferStatus(tx, db.EthTransferdoneTb{
				Txhash: transfer.Txhash,
				State:  int64(status),
				Msg:    errMsg,
			})
			if err != nil {
				return err
			}

			// 交易状态不管完成与否,都不再查询,从表里移除
			err = db.DeleteTransfer(tx, transfer.ID)

			return err
		})

		log.Infof("checkTransfer Transaction err : %v", err)
	}
}

func transactionReceipt(c *erc20.ReturnClient, hash string) (uint64, error) {
	srt, err := c.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		log.Errorf("transactionReceipt err : %v", err)
		return 0, err
	}
	// log.Infof("srt------------Status:%v,BlockNumber:%v", srt.Status, srt.BlockNumber)

	return srt.Status, nil
}

// getBlockNumber 获取区块高度
func getBlockNumber(c *erc20.ReturnClient) (uint64, error) {
	bn, err := c.Client.BlockNumber(context.Background())
	if err != nil {
		log.Errorf("getBlockNumber BlockNumber err : %v", err)
	}

	return bn, err
}
