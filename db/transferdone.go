package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// eth_transferdone_key 表

// SaveTransferStatus 保存消息状态
func SaveTransferStatus(txhash string, status uint64, msg string) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	execStr := fmt.Sprintf("insert into %s(txhash,status,msg) values(?,?,?)", transferDoneTable)

	_, err := db.Exec(execStr, txhash, status, msg)
	if err != nil {
		log.Errorln("sql SaveTransferStatus err : ", err)
	}

	return err
}
