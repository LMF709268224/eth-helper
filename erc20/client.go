package erc20

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
	Filter     *TokenERC20Filterer
	Caller     *TokenERC20Caller
	ReCaller   *TokenRecipientCaller
}

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
