package ethevent

import (
	"eth-helper/erc20"

	log "github.com/sirupsen/logrus"
)

var messageChan = make(chan *erc20.PkgnameTransfer, 100)

// Init init watch
func Init() {
	go watchChan()

	watchTransfer()

	// test
	// testChan()

	// go test()
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

	// for {
	// 	select {
	// 	case a := <-messageChan:
	// 		log.Infof("from=%s,to=%s,数量=%d,交易hash=%s", a.From, a.To, a.Value.Int64(), a.Raw.TxHash)
	// 	}
	// }
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
