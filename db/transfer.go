package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// eth_transfer_tbs 表

// SaveNewTransfer 保存待确认消息
func SaveNewTransfer(info EthTransferTb) error {
	db := getDBConnection()

	tx := db.Create(&info)
	if tx.Error != nil {
		log.Errorf("SaveNewTransfer Error :%v", tx.Error)
	}

	return tx.Error
}

// GetMinBlocknumber 获取最低blocknumber
func GetMinBlocknumber() (uint64, error) {
	db := getDBConnection()

	blocknumber := uint64(0)

	ss := fmt.Sprintf("select MIN(blocknumber) from  %s", transferTable)
	row := db.Raw(ss)

	tx := row.Scan(&blocknumber)
	if tx.Error != nil && tx.Error.Error() != "sql: Scan error on column index 0" {
		log.Errorf("GetMinBlocknumber err : %v", tx.Error)
	}
	// fmt.Printf("blocknumber : %v", blocknumber)

	return blocknumber, tx.Error
}

// GetTransferWithBlocknumber 获取某一高度的全部交易
func GetTransferWithBlocknumber(blocknumber uint64) ([]EthTransferTb, error) {
	db := getDBConnection()

	var rows []EthTransferTb
	// select
	tx := db.Where("blocknumber = ?", blocknumber).Find(&rows)
	if tx.Error != nil {
		log.Errorf("GetTransferWithBlocknumber Error :%v", tx.Error)
	}

	return rows, tx.Error
}

// DeleteTransfer 删除某个交易
func DeleteTransfer(mid int64) error {
	db := getDBConnection()

	var row []EthTransferTb

	tx := db.Where("id = ?", mid).Delete(&row)
	if tx.Error != nil {
		log.Errorf("DeleteTransfer Error :%v", tx.Error)
	}

	return tx.Error
}

// // SaveNewTransfer 保存待确认消息
// func SaveNewTransfer(info EthTransferTb) error {
// 	db := openDB()
// 	// 延迟到函数结束关闭链接
// 	defer db.Close()

// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Errorln("sql SaveNewTransfer Begin err : ", err)
// 		return err
// 	}

// 	execStr := fmt.Sprintf("insert into %s(mto,mfrom,txhash,value,blocknumber) values(?,?,?,?,?)", transferTable)

// 	_, err = tx.Exec(execStr, info.MTo, info.MFrom, info.Txhash, info.Value, info.Blocknumber)
// 	if err != nil {
// 		log.Errorln("sql SaveNewTransfer Exec err : ", err)
// 		return err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		fmt.Println("sql SaveNewTransfer Commit err :", err)
// 	}

// 	return err
// }

// // GetMinBlocknumber 获取最低blocknumber
// func GetMinBlocknumber() (blocknumber uint64, err error) {
// 	db := openDB()
// 	// 延迟到函数结束关闭链接
// 	defer db.Close()

// 	// 执行单条查询
// 	ss := fmt.Sprintf("select MIN(blocknumber) from  %s", transferTable)
// 	// ss := fmt.Sprintf("select blocknumber from %s limit 1", transferTable)
// 	row := db.QueryRow(ss)

// 	err = row.Scan(&blocknumber)
// 	if err != nil {
// 		if err.Error() == "sql: Scan error on column index 0" {
// 			log.Errorln("hdgfjhasgdfjkhsdjhfsjahfjs jfsagfjh ")
// 		}
// 		log.Errorf("GetMinBlocknumber Scan err : %v", err)
// 	}

// 	return
// }

// // GetTransferWithBlocknumber 获取某一高度的全部交易
// func GetTransferWithBlocknumber(blocknumber uint64) ([]EthTransferTb, error) {
// 	db := openDB()
// 	// 延迟到函数结束关闭链接
// 	defer db.Close()

// 	ss := fmt.Sprintf("select * from %s where blocknumber = ?", transferTable)
// 	rows, err := db.Query(ss, blocknumber)

// 	defer func() {
// 		if rows != nil {
// 			rows.Close() // 可以关闭掉未scan连接一直占用
// 		}
// 	}()

// 	if err != nil {
// 		log.Errorf("GetTransferWithBlocknumber Query err:%v", err)
// 		return nil, err
// 	}

// 	infos := []EthTransferTb{}

// 	for rows.Next() {
// 		info := EthTransferTb{}

// 		err = rows.Scan(&info.ID, &info.MTo, &info.MFrom, &info.Txhash, &info.Value, &info.Blocknumber) // 不scan会导致连接不释放
// 		if err != nil {
// 			log.Errorf("GetTransferWithBlocknumber Scan err:%v", err)
// 			return nil, err
// 		}

// 		log.Infof("GetTransferWithBlocknumber info : %v", info)

// 		infos = append(infos, info)
// 	}

// 	return infos, nil
// }

// // DeleteTransfer 删除某个交易
// func DeleteTransfer(mid int64) error {
// 	db := openDB()
// 	// 延迟到函数结束关闭链接
// 	defer db.Close()

// 	ss := fmt.Sprintf("delete from %s where id = ?", transferTable)

// 	_, err := db.Exec(ss, mid)
// 	if err != nil {
// 		fmt.Printf("DeleteTransfer err:%v", err)
// 	}

// 	return err
// }
