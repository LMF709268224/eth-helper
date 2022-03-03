package ethevent

import (
	"context"
	"eth-helper/db"
	"eth-helper/erc20"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// var messageChan = make(chan *erc20.TokenERC20Transfer, 100)

// Init init watch
func Init() {
	// go watchChan()

	go watchTransfer()

	// test
	// testChan()

	// go test()
}

// func watchChan() {
// 	for {
// 		select {
// 		case a := <-messageChan:
// 			log.Infof("from=%s,to=%s,数量=%d,交易hash=%s", a.From, a.To, a.Value.Int64(), a.Raw.TxHash)
// 		}
// 	}
// }

//  监听盲盒价格变化
func watchTransfer() {
	c := erc20.GetClient()
	defer c.Client.Close()

	messageChan := make(chan *erc20.TokenERC20Transfer, 100)
	// from := []common.Address{}
	// to := []common.Address{common.HexToAddress("0x6A859c1BD1D0722D7Bc3Df05a96FF7684dCA30eD")}
	sub, err := c.Filter.WatchTransfer(nil, messageChan, nil, nil)
	if err != nil {
		log.Errorf("WatchTransfer err : %v", err)
		return
	}

	log.Infoln("WatchTransfer------------")

	srt, err := c.Client.TransactionReceipt(context.Background(), common.HexToHash("0xb14b073a8e87aa0c82d254c363dd05a04cee2eb705be704818d37f7adb82b4f3"))

	log.Infof("srt------------Status:%v,BlockNumber:%v", srt.Status, srt.BlockNumber)
	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				log.Errorf("eWatchTransfer transfer err:%v", err)
			}
		case transfer := <-messageChan:
			// log.Infof("from=%s,to=%s,数量=%d,交易hash=%s", transfer.From, transfer.To, transfer.Value.Int64(), transfer.Raw.TxHash)
			log.Infof("from=%s,to=%s,数量=%d,交易hash=%s,BlockNumber=%v", transfer.From.Hex(), transfer.To.Hex(), transfer.Value.Int64(), transfer.Raw.TxHash.Hex(), transfer.Raw.BlockNumber)
			// newTransfer(transfer)
		}
	}
}

func newTransfer(transfer *erc20.TokenERC20Transfer) {
	db.SaveNewTransfer(transfer.To.Hex(), transfer.From.Hex(), transfer.Raw.TxHash.Hex(), transfer.Value.Int64())
}

// func testChan() {
// 	// test
// 	messageChan <- &erc20.PkgnameTransfer{
// 		From:  common.HexToAddress("123456"),
// 		To:    common.HexToAddress("00000000000000000"),
// 		Value: big.NewInt(int64(256788)),
// 		Raw:   types.Log{},
// 	}
// }

// func test() {
// 	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/a5c713d632f944df9a77d56cf08f9083")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
// 	query := ethereum.FilterQuery{
// 		Addresses: []common.Address{contractAddress},
// 	}

// 	logs := make(chan types.Log)

// 	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			log.Fatal(err)
// 		case vLog := <-logs:
// 			fmt.Println(vLog) // pointer to event log
// 			log.Infoln("WatchTransfer------------")
// 		}
// 	}
// }
