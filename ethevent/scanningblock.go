package ethevent

import (
	"bytes"
	"context"
	"encoding/json"
	"eth-helper/db"
	"eth-helper/erc20"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	blocknumber int

	myAddress map[string]db.EthAddressTb

	chainid = 4

	contractAddress = ""
)

// InitScanningBlockTask 初始化扫快
func InitScanningBlockTask(num int, ca string) {
	// 合约地址
	contractAddress = ca
	// 获取所有地址
	myAddress = db.GetAllAddress()
	for k, _ := range myAddress {
		log.Println("a:", k)
	}

	blocknumber = num

	timeSecond := 60 * 1
	duration := time.Duration(time.Second * time.Duration(timeSecond))

	t := time.NewTicker(duration)
	defer t.Stop()

	for {
		<-t.C

		log.Infoln("ScanningBlock task...:")
		scanningBlock()
	}
}

func scanningBlock() {
	c := erc20.GetClient()
	defer c.Client.Close()

	// 读数据库看看数据库里的高度 blocknumberDB + 1
	blocknumberDB := db.GetBlockNumber() + 1

	// 对比配置高度 用值大的来开始扫
	if blocknumberDB > blocknumber {
		blocknumber = blocknumberDB
	}

	log.Infoln("blocknumber...:", blocknumber)
	err := workerHander(blocknumber)
	if err != nil {
		return
	}
	// // 扫快
	// block, err := getBlockByNumber(c, int64(blocknumber))
	// if err != nil {
	// 	return
	// }

	// // 处理订单
	// err = readTransactions(block)
	// if err != nil {
	// 	return
	// }

	// 最新num存DB
	db.SaveBlockNumber(blocknumber)

	// blocknumber ++
	blocknumber++
}

// 具体处理区块
func workerHander(num int) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
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

	for i := 0; i < len(txs); i++ {
		tx := txs[i].(map[string]interface{})
		if tx["to"] == nil { // 部署合约交易直接跳过
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
		// log.Info("find_a_token")

		// 验证合约
		contract := tx["to"].(string)

		log.Info("contract : ", contract)
		// exists, _, deci := utils.IsExistsContract(contract)
		// if !exists {
		// 	log.Info("合约不一致：", contract)
		// 	continue
		// }
		if strings.EqualFold(contract, contractAddress) == false {
			log.Info("合约不一致：", contract)
			continue
		}

		// fee_address := conf.Cfg.MustValue("eth", "fee_address")
		// if fmt.Sprintf("0x%s", input[34:74]) == fee_address {
		// 	continue
		// }
		// out_from_address := conf.Cfg.MustValue("eth", "out_from_address")
		// if fmt.Sprintf("0x%s", input[34:74]) == out_from_address {
		// 	continue
		// }

		ext := existsAddress(fmt.Sprintf("0x%s", input[34:74]), chainid, tx["to"].(string))
		// log.Info("测试erc20:", fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string))
		if !ext {
			continue
		}

		log.Info("地址存在：", fmt.Sprintf("0x%s", input[34:74]), chainid, tx["to"].(string))

		// 检查erc20代币转账是否成功
		// if p.checkERC20IsSuccess(p.Url, tx["hash"].(string)) == false {
		// 	log.Info("checkERC20IsSuccess_continue", tx["hash"].(string))
		// 	continue
		// }

		// walletToken := new(models.WalletToken)
		// walletToken.GetUidByAddress(fmt.Sprintf("0x%s", input[34:74]))

		// log.Info("用户uid：", walletToken.Uid)

		// ok, err := p.newOrder(walletToken.Uid, tx["from"].(string), fmt.Sprintf("0x%s", input[34:74]), p.Chainid, tx["to"].(string), fmt.Sprintf("0x%s", input[vstart:138]), tx["hash"].(string), tx["gas"].(string), tx["gasPrice"].(string), deci)
		// log.Info("newOrder complete", ok, err)
		// continue

	}
	return nil
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
	if _, ok := myAddress[field]; !ok {
		// 不是我们自己的地址
		return false
	}
	return true
}

// post RPC数据
func post(send map[string]interface{}) ([]byte, error) {
	bytesData, err := json.Marshal(send)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	url := "https://rinkeby.infura.io/v3/a5c713d632f944df9a77d56cf08f9083"

	log.Info(url)

	request, err := http.NewRequest("POST", url, reader)
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

			if txHash == "0xeb377059a7814e49ba1cbd4800abf8b7aa0c99f648e1fd2b0ef32785a4750bbd" {
				log.Infoln("readTransactions find to :", to)
			}

			tolower := strings.ToLower(to)
			// log.Infoln("readTransactions find tolower :", tolower)
			if _, ok := myAddress[tolower]; !ok {
				// 不是我们自己的地址
				continue
			}

			log.Infoln("readTransactions txHash :", txHash)

			// 看看订单是否已经处理 eth_transferdone_tbs
			_, err := db.GetTransferInfo(tx, txHash)
			if err == nil {
				// 在表里找到了 则不处理
				log.Infoln("readTransactions在表里找到了 txHash :", txHash)
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
