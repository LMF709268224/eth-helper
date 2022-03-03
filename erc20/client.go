package erc20

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
