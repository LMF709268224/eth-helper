package message

import (
	"eth-helper/erc20"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

var messageChan = make(chan *erc20.PkgnameTransfer, 100)

// Init init watch
func Init() {
	go watchChan()

	watchTransfer()

	// test
	testChan()
}

func watchChan() {
	for {
		select {
		case a := <-messageChan:
			log.Infof("from=%s,to=%s,数量=%d,交易hash=%s", a.From, a.To, a.Value.Int64(), a.Raw.TxHash)
		}
	}
}

//  监听盲盒价格变化
func watchTransfer() {
	c := erc20.GetClient()
	defer c.Client.Close()

	// from := []common.Address{}
	// to := []common.Address{}
	_, err := c.Filter.WatchTransfer(nil, messageChan, nil, nil)
	if err != nil {
		log.Errorf("WatchTransfer err : %v", err)
		return
	}

	log.Infoln("WatchTransfer------------")
}

func testChan() {
	// test
	messageChan <- &erc20.PkgnameTransfer{
		From:  common.HexToAddress("123456"),
		To:    common.HexToAddress("00000000000000000"),
		Value: big.NewInt(int64(256788)),
		Raw:   types.Log{},
	}
}
