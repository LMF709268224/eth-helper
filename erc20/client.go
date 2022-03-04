package erc20

import (
	"crypto/ecdsa"
	"errors"
	"eth-helper/db"
	"fmt"
	"math/big"

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

// GetClient 获取ethclient
func GetClient() *ReturnClient {
	node := "wss://rinkeby.infura.io/ws/v3/4a500de3b58c4ee29f06f412c041669c"
	client, err := ethclient.Dial(node)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	contract := "0xF71B99E8c9EF7fe986C9Ff3A4913855854f28C4D"

	c, err := NewTokenERC20Filterer(common.HexToAddress(contract), client)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	caller, err := NewTokenERC20Caller(common.HexToAddress(contract), client)
	if err != nil {
		return nil
	}

	Recaller, err := NewTokenRecipientCaller(common.HexToAddress(contract), client)
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
func NewAddress() (db.MAddressInfo, error) {
	info := db.MAddressInfo{}

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
	// info.Pwd = privateKeyBytes TODO 私钥加密

	return info, nil
}
