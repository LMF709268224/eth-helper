package message

import (
	"context"
	"eth-helper/erc20"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// TestHttp test
func TestHttp() {
	c := erc20.GetHttpClient()
	defer c.Client.Close()

	i, err := c.Client.ChainID(context.Background())
	if err != nil {
		log.Errorf("ChainID err : %v", err)
	}
	log.Infof("ChainID i : %v", i)

	b, err := c.Client.BlockNumber(context.Background())
	if err != nil {
		log.Errorf("BlockNumber err : %v", err)
	}
	log.Infof("BlockNumber b : %v", b)

	add := "0x7ee5ecee43bd050238d08b56d933398caacb5044"
	address := common.HexToAddress(add)
	balance, err := c.Client.BalanceAt(context.Background(), address, big.NewInt(int64(b)))
	if err != nil {
		log.Errorf("BalanceAt err : %v", err)
	}
	log.Infof("BalanceAt balance : %v", balance)

	// c.Client.get
}
