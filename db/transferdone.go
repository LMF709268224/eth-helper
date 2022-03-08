package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// eth_transferdone_tbs 表

// SaveTransferStatus 保存消息状态
func SaveTransferStatus(db *gorm.DB, info EthTransferdoneTb) error {
	// if db == nil {
	// 	db = GetDBConnection()
	// }

	tx := db.Create(&info)
	if tx.Error != nil {
		log.Errorf("SaveNewAddress Error :%v", tx.Error)
	}

	return tx.Error
}

// // SaveTransferStatus 保存消息状态
// func SaveTransferStatus(txhash string, status uint64, msg string) error {
// 	db := openDB()
// 	// 延迟到函数结束关闭链接
// 	defer db.Close()

// 	execStr := fmt.Sprintf("insert into %s(txhash,status,msg) values(?,?,?)", transferDoneTable)

// 	_, err := db.Exec(execStr, txhash, status, msg)
// 	if err != nil {
// 		log.Errorln("sql SaveTransferStatus err : ", err)
// 	}

// 	return err
// }
