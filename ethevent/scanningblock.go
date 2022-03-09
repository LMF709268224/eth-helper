package ethevent

import (
	"bytes"
	"encoding/json"
	"eth-helper/db"
	"eth-helper/erc20"
	"eth-helper/redishelper"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ouqiang/timewheel"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	blocknumber         int
	chainid             = 4
	contract            = ""
	nodehttps           = ""
	ethCheckCBTranNewTW *timewheel.TimeWheel // 时间轮，检查交易状态
)

// InitScanningBlockTask 初始化扫快
func InitScanningBlockTask(blockNumber int, contractAddress string, chainID int, nodeHTTPS string) {
	// defer func() {
	// 	// 删除定时器, 参数为添加定时器传递的唯一标识
	// 	ethCheckCBTranNewTW.RemoveTimer("eth_check_cbi")
	// 	// 停止时间轮
	// 	ethCheckCBTranNewTW.Stop()
	// }()
	nodehttps = nodeHTTPS
	chainid = chainID
	// 合约地址
	contract = contractAddress

	blocknumber = blockNumber

	ethCheckCBTranNewTW = timewheel.New(1*time.Second, 3600, func(data interface{}) {
		fmt.Println("start eth.cbi.watch.new...")
		// 区块操作处理
		scanningBlock()
		// 继续添加定时器
		ethCheckCBTranNewTW.AddTimer(5*time.Second, "eth_check_cbi", nil)
	})
	ethCheckCBTranNewTW.Start()
	// 开始一个事件处理
	ethCheckCBTranNewTW.AddTimer(5*time.Second, "eth_check_cbi", nil)
}

func scanningBlock() {
	c := erc20.GetClient()
	defer c.Client.Close()

	// 读数据库看看数据库里的高度 blocknumberDB + 1
	blocknumberDB := redishelper.GetBlockNumber(chainid) + 1

	// 对比配置高度 用值大的来开始扫
	if blocknumberDB > blocknumber {
		blocknumber = blocknumberDB
	}

	// log.Infoln("blocknumber...:", blocknumber)
	// 扫快  处理订单
	err := workerHander(blocknumber)
	if err != nil {
		log.Infoln("workerHander err :", err.Error())
		return
	}

	// 最新num存DB
	err = redishelper.SaveBlockNumber(chainid, blocknumber)
	if err != nil {
		log.Infoln("SaveBlockNumber err :", err.Error())
	}

	// blocknumber ++
	blocknumber++
}

// 具体处理区块
func workerHander(num int) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误:%s\n", r)
		}
	}()

	ret, err := getblockBynumber(num)
	if err != nil {
		log.Error(err)
		return err
	}
	var block map[string]interface{}

	json.Unmarshal(ret, &block)
	txs := block["result"].(map[string]interface{})["transactions"].([]interface{})

	// log.Info(txs)

	// 在事务中处理
	mdb := db.GetDBConnection()
	err = mdb.Transaction(func(gdb *gorm.DB) error {
		for i := 0; i < len(txs); i++ {
			tx := txs[i].(map[string]interface{})
			if tx["to"] == nil { // 不属合约交易直接跳过
				continue
			}

			input := tx["input"].(string)
			// 不是token转账跳过
			if strings.Count(input, "") < 138 || strings.Compare(input[0:10], "0xa9059cbb") != 0 {
				continue
			}

			var vstart int
			for i := 74; i < 138; i++ {
				if input[i:i+1] != "0" {
					vstart = i
					break
				}
			}
			if vstart == 0 {
				continue
			}

			// 验证合约
			cTo := tx["to"].(string)
			// log.Info("contract : ", contract)
			if strings.EqualFold(cTo, contract) == false {
				continue
			}

			toAddress := fmt.Sprintf("0x%s", input[34:74])
			ext := existsAddress(toAddress, chainid, tx["to"].(string))
			if !ext {
				continue
			}

			txhash := tx["hash"].(string)
			from := tx["from"].(string)

			value := fmt.Sprintf("0x%s", input[vstart:138])
			temp, _ := new(big.Int).SetString(value[2:], 16)
			// float64_value, _ := decimal.NewFromBigInt(temp, int32(-deci)).Float64()

			err = saveAndDelDb(gdb, txhash, toAddress, temp.String(), num, from)
			if err != nil {
				return err
			}
		}

		return nil
	})

	log.Printf("Transaction err : %v", err)

	return nil
}

func saveAndDelDb(tx *gorm.DB, txHash string, to string, value string, num int, from string) error {
	// 看看订单是否已经处理 eth_transferdone_tbs
	info, err := db.GetTransferInfo(tx, txHash)
	if err == nil && info.Txhash == txHash {
		// 在表里找到了 则不处理
		log.Infoln("saveAndDelDb 在表里找到了 txHash :", txHash)
		return nil
	}

	// 写到待处理表 eth_transfer_tbs
	err = db.SaveNewTransfer(tx, db.EthTransferTb{
		MTo:         to,
		Txhash:      txHash,
		Value:       value,
		Blocknumber: uint64(num),
		MFrom:       from,
	})
	if err != nil {
		log.Errorf("saveAndDelDb SaveNewTransfer err : %v,hash : %s", err, txHash)
	}

	return err
}

func getblockBynumber(num int) ([]byte, error) {
	send := make(map[string]interface{})
	send["jsonrpc"] = "2.0"
	send["method"] = "eth_getBlockByNumber"
	strconv.FormatInt(int64(num), 16)
	send["params"] = []interface{}{fmt.Sprintf("0x%s", strconv.FormatInt(int64(num), 16)), true}
	send["id"] = chainid // https://learnblockchain.cn/article/1791
	rsp, err := post(send)

	return rsp, err
}

// 判断地址是否存在
func existsAddress(address string, chainid int, contract string) bool {
	// key := "ETH_CBI_ADDRESS_REDIS_KEY"
	// field := "%s" //chainid:address:contract
	// field = fmt.Sprintf(field, address)
	// field = strings.ToLower(field)
	// return utils.Redis.HExists(key, field).Val()
	field := strings.ToLower(address)

	return redishelper.ExistsAddress(field)
}

// post RPC数据
func post(send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	reader := bytes.NewReader(bytesData)

	// log.Info(url)

	request, err := http.NewRequest("POST", nodehttps, reader)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// byte数组直接转成string，优化内存
	return respBytes, nil
}

// func readTransactions(block *types.Block) error {
// 	// 在事务中处理
// 	mdb := db.GetDBConnection()
// 	err := mdb.Transaction(func(tx *gorm.DB) error {
// 		for _, transaction := range block.Transactions() {

// 			if transaction.To() == nil {
// 				continue
// 			}

// 			to := transaction.To().Hex()
// 			txHash := transaction.Hash().Hex()

// 			if txHash == "0xeb377059a7814e49ba1cbd4800abf8b7aa0c99f648e1fd2b0ef32785a4750bbd" {
// 				log.Infoln("readTransactions find to :", to)
// 			}

// 			tolower := strings.ToLower(to)
// 			// log.Infoln("readTransactions find tolower :", tolower)
// 			if _, ok := myAddress[tolower]; !ok {
// 				// 不是我们自己的地址
// 				continue
// 			}

// 			log.Infoln("readTransactions txHash :", txHash)

// 			// 看看订单是否已经处理 eth_transferdone_tbs
// 			_, err := db.GetTransferInfo(tx, txHash)
// 			if err == nil {
// 				// 在表里找到了 则不处理
// 				log.Infoln("readTransactions在表里找到了 txHash :", txHash)
// 				continue
// 			}

// 			// 写到待处理表 eth_transfer_tbs
// 			err = db.SaveNewTransfer(tx, db.EthTransferTb{
// 				MTo:    to,
// 				Txhash: txHash,
// 				Value:  transaction.Value().String(),
// 			})
// 			if err != nil {
// 				log.Errorf("readTransactions SaveNewTransfer err : %v,hash : %s", err, txHash)
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	return err
// }

// func getBlockByNumber(c *erc20.ReturnClient, num int64) (*types.Block, error) {
// 	block, err := c.Client.BlockByNumber(context.Background(), big.NewInt(num))
// 	if err != nil {
// 		log.Errorf("getBlockByNumber BlockByNumber err : %v", err)
// 	}

// 	return block, err
// }

// // TestGetBlock 测试
// func TestGetBlock() {
// 	c := erc20.GetClient()
// 	defer c.Client.Close()

// 	num, err := getBlockNumber(c)
// 	if err != nil {
// 		return
// 	}
// 	log.Infoln("getBlockNumber num :", num)

// 	block, err := getBlockByNumber(c, int64(num))
// 	if err != nil {
// 		return
// 	}
// 	log.Infof("getBlockByNumber : %v", block)

// 	for _, tx := range block.Transactions() {
// 		// log.Infof("tx:%v", tx)
// 		// log.Infof("Hash: %s ,Value:%s, To:%s", tx.Hash().Hex(), tx.Value().String(), tx.To().Hex())
// 		log.Infoln(tx.Hash().Hex())     // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
// 		log.Infoln(tx.Value().String()) // 10000000000000000
// 		log.Infoln(tx.To())             // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
// 		log.Infoln("---------------------")
// 	}
// }

// func getBlockHeaderByNumber(c *erc20.ReturnClient, num int64) (*types.Header, error) {
// 	header, err := c.Client.HeaderByNumber(context.Background(), big.NewInt(num))
// 	if err != nil {
// 		log.Errorf("getBlockHeaderByNumber HeaderByNumber err : %v", err)
// 	}

// 	return header, err
// }

// func getBlockByHash(c *erc20.ReturnClient, hash common.Hash) (*types.Block, error) {
// 	block, err := c.Client.BlockByHash(context.Background(), hash)
// 	if err != nil {
// 		log.Errorf("getBlockByHash BlockByHash err : %v", err)
// 	}

// 	return block, err
// }
