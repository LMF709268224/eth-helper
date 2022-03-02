package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ReturnClient struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	GasPrice   *big.Int
	Nonce      int64
	Filter     *PkgnameFilterer
}

func GetClient() *ReturnClient {
	node := "节点地址 ws"
	client, err := ethclient.Dial(node)
	if err != nil {
		fmt.Println("Failed to Dial ", err)
		return nil
	}

	contract := "0xdAC17F958D2ee523a2206206994597C13D831ec7"

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
