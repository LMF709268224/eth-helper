package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	go WatchTransfer()
	select {}
}

//监听盲盒价格变化
func WatchTransfer() {
	c := GetClient()
	defer c.Client.Close()

	ch := make(chan *PkgnameTransfer, 100)
	from := []common.Address{}
	to := []common.Address{}
	c.Filter.WatchTransfer(nil, ch, from, to)
	for {
		select {
		case a := <-ch:
			fmt.Sprintf("from=%s,to=%s,数量=%d,交易hash=%s", a.From, a.To, a.Value.Int64(), a.Raw.TxHash)
		}
	}
}
