package ethevent

import (
	"eth-helper/db"
	"eth-helper/erc20"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// InitWatchTransfer 初始化监听交易
func InitWatchTransfer() {
	// go watchChan()
	list := db.GetAllAddressInfo()
	if len(list) <= 0 {
		log.Errorln("InitWatchTransfer list is nil...")
		return
	}

	// 监听地址
	to := []common.Address{}
	for _, info := range list {
		to = append(to, common.HexToAddress(info.Address))
	}

	go watchTransfer(to)
}

//  监听盲盒价格变化
func watchTransfer(to []common.Address) {
	c := erc20.GetClient()
	defer c.Client.Close()

	messageChan := make(chan *erc20.TokenERC20Transfer, 100)
	// from := []common.Address{}
	// to := []common.Address{common.HexToAddress("0x6A859c1BD1D0722D7Bc3Df05a96FF7684dCA30eD")}
	sub, err := c.Filter.WatchTransfer(nil, messageChan, nil, to)
	if err != nil {
		log.Errorf("WatchTransfer err : %v", err)
		return
	}

	log.Infoln("WatchTransfer------------")

	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				log.Errorf("watchTransfer transfer err:%v", err)
			}
		case transfer := <-messageChan:
			// log.Infof("from=%s,to=%s,数量=%d,交易hash=%s", transfer.From, transfer.To, transfer.Value.Int64(), transfer.Raw.TxHash)
			log.Infof("from=%s,to=%s,数量=%d,交易hash=%s", transfer.From.Hex(), transfer.To.Hex(), transfer.Value.Int64(), transfer.Raw.TxHash.Hex())
			err := newTransfer(transfer)
			if err != nil {
				log.Errorf("watchTransfer newTransfer err : %v,hash : %s", err, transfer.From.Hex())
			}
		}
	}
}

func newTransfer(transfer *erc20.TokenERC20Transfer) error {
	info := db.MTransferInfo{
		To:          transfer.To.Hex(),
		From:        transfer.From.Hex(),
		Txhash:      transfer.Raw.TxHash.Hex(),
		Value:       transfer.Value.Int64(),
		Blocknumber: transfer.Raw.BlockNumber,
	}

	return db.SaveNewTransfer(info)
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