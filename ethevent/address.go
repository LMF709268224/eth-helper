package ethevent

import (
	"eth-helper/db"
	"eth-helper/erc20"

	log "github.com/sirupsen/logrus"
)

// NewAddresss num:address num
func NewAddresss(num int) error {
	infos := []db.EthAddressTb{}

	for i := 0; i < num; i++ {
		info, err := erc20.NewAddress()
		if err != nil {
			log.Errorf("NewAddresss NewAddress err:%v", err)
			continue
		}

		infos = append(infos, info)
	}

	// save db
	err := db.SaveNewAddress(infos)
	if err != nil {
		log.Errorf("NewAddresss SaveNewAddress err:%v", err)
	}

	return err
}
