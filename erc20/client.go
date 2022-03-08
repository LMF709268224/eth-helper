package erc20

import (
	"crypto/ecdsa"
	"errors"
	"eth-helper/db"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

// ReturnClient ethclient 结构体
type ReturnClient struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	GasPrice   *big.Int
	Nonce      int64
	Filter     *TokenERC20Filterer
	Caller     *TokenERC20Caller
	ReCaller   *TokenRecipientCaller
}

var (
	nodewss         = ""
	contractaddress = ""
)

// Init 初始化 节点wss 合约地址
func Init(ns, ca string) {
	nodewss = ns
	contractaddress = ca
}

// GetClient 获取ethclient
func GetClient() *ReturnClient {
	// node := "wss://rinkeby.infura.io/ws/v3/4a500de3b58c4ee29f06f412c041669c"
	client, err := ethclient.Dial(nodewss)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	// contract := "0xF71B99E8c9EF7fe986C9Ff3A4913855854f28C4D"

	c, err := NewTokenERC20Filterer(common.HexToAddress(contractaddress), client)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	caller, err := NewTokenERC20Caller(common.HexToAddress(contractaddress), client)
	if err != nil {
		return nil
	}

	Recaller, err := NewTokenRecipientCaller(common.HexToAddress(contractaddress), client)
	if err != nil {
		return nil
	}

	return &ReturnClient{
		Client:   client,
		Filter:   c,
		Caller:   caller,
		ReCaller: Recaller,
	}
}

// NewAddress 生成地址
func NewAddress() (db.EthAddressTb, error) {
	info := db.EthAddressTb{}

	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Errorf("newAddress GenerateKey err :%v", err)
		return info, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("私钥为: " + hexutil.Encode(privateKeyBytes))

	// 私钥导出地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Errorln("newAddress cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return info, errors.New("publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("地址为: " + address)

	info.Address = address
	info.Pwd = hexutil.Encode(privateKeyBytes)
	info.AddTime = time.Now()
	// info.Pwd = privateKeyBytes TODO 私钥加密

	return info, nil
}

// func main() {
// 	privKey, address, err := KeystoreToPrivateKey("UTC--2017-11-21T05-46-23.555205600Z--6e60f5243e1a3f0be3f407b5afe9e5395ee82aa2", "123456789")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("privKey:%s\naddress:%s\n", privKey, address)
// }

// func KeystoreToPrivateKey(privateKeyFile, password string) (string, string, error) {
// 	keyjson, err := ioutil.ReadFile(privateKeyFile)
// 	if err != nil {
// 		fmt.Println("read keyjson file failed：", err)
// 	}
// 	unlockedKey, err := keystore.DecryptKey(keyjson, password)
// 	if err != nil {

// 		return "", "", err

// 	}
// 	privKey := hex.EncodeToString(unlockedKey.PrivateKey.D.Bytes())
// 	addr := crypto.PubkeyToAddress(unlockedKey.PrivateKey.PublicKey)
// 	return privKey, addr.String(), nil
// }

// func createKs() {
// 	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
// 	password := "123"
// 	account, err := ks.NewAccount(password)
// 	if err != nil {
// 		log.Info("er:", err)
// 	}

// 	fmt.Println(account.Address.Hex())
// }

// func importKs() {
// 	file := "./tmp/UTC--2022-03-03T04-30-43.877689000Z--ee01812ad65ca9ebd3ea9619aa007d73326b7570"
// 	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
// 	jsonBytes, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		log.Info("er:", err)
// 	}

// 	password := "123"
// 	account, err := ks.Import(jsonBytes, password, password)
// 	if err != nil {
// 		log.Info("er:", err)
// 	}

// 	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3

// 	if err := os.Remove(file); err != nil {
// 		log.Info("er:", err)
// 	}
// }
