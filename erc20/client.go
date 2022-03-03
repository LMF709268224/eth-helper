package erc20

import (
	"crypto/ecdsa"
	"errors"
	"eth-helper/db"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// ReturnClient Client struct
type ReturnClient struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	GasPrice   *big.Int
	Nonce      int64
	Filter     *PkgnameFilterer
}

// GetClient get the Client
func GetClient() *ReturnClient {
	// test := "wss://rinkeby.infura.io/ws/v3/a5c713d632f944df9a77d56cf08f9083"
	node := "wss://mainnet.infura.io/ws/v3/a5c713d632f944df9a77d56cf08f9083"
	client, err := ethclient.Dial(node)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	// testContract := "0x1B93Ae4E96fDEA0B49ef51F77B388b45160dEd90"
	contract := "0xdAC17F958D2ee523a2206206994597C13D831ec7" // usdt
	// contract := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // usdc

	c, err := NewPkgnameFilterer(common.HexToAddress(contract), client)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	return &ReturnClient{
		Client: client,
		Filter: c,
	}
}

// GetHttpClient http Client
func GetHttpClient() *ReturnClient {
	node := "https://mainnet.infura.io/v3/a5c713d632f944df9a77d56cf08f9083"
	client, err := ethclient.Dial(node)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	contract := "0xdAC17F958D2ee523a2206206994597C13D831ec7" // usdt
	// contract := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // usdc

	c, err := NewPkgnameFilterer(common.HexToAddress(contract), client)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	return &ReturnClient{
		Client: client,
		Filter: c,
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
	// info.PrivateKey = privateKeyBytes
	// info.Pwd = privateKeyBytes TODO 私钥加密

	return info, nil
}
