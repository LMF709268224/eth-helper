package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// eth_transfer_key 表

// SaveNewTransfer 保存待确认消息
func SaveNewTransfer(info MTransferInfo) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	execStr := fmt.Sprintf("insert into %s(mto,mfrom,txhash,value,blocknumber) values(?,?,?,?,?)", transferTable)

	_, err := db.Exec(execStr, info.To, info.From, info.Txhash, info.Value, info.Blocknumber)
	if err != nil {
		log.Errorln("sql SaveNewTransfer err : ", err)
	}

	return err
}

// GetFristTransfer 获取第一个,看看blocknumber
func GetFristTransfer() (blocknumber uint64, err error) {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	// 执行单条查询
	ss := fmt.Sprintf("select blocknumber from %s limit 1", transferTable)
	row := db.QueryRow(ss)

	err = row.Scan(&blocknumber)
	if err != nil {
		log.Errorf("GetFristTransfer Scan err : %v", err)
	}

	return
}

// GetTransferWithBlocknumber 获取某一高度的全部交易
func GetTransferWithBlocknumber(blocknumber uint64) ([]MTransferInfo, error) {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	ss := fmt.Sprintf("select * from %s where blocknumber = ?", transferTable)
	rows, err := db.Query(ss, blocknumber)

	defer func() {
		if rows != nil {
			rows.Close() // 可以关闭掉未scan连接一直占用
		}
	}()

	if err != nil {
		log.Errorf("GetTransferWithBlocknumber Query err:%v", err)
		return nil, err
	}

	infos := []MTransferInfo{}

	for rows.Next() {
		info := MTransferInfo{}

		err = rows.Scan(&info.ID, &info.To, &info.From, &info.Txhash, &info.Value, &info.Blocknumber) // 不scan会导致连接不释放
		if err != nil {
			log.Errorf("GetTransferWithBlocknumber Scan err:%v", err)
			return nil, err
		}

		log.Infof("GetTransferWithBlocknumber info : %v", info)

		infos = append(infos, info)
	}

	return infos, nil
}

// DeleteTransfer 删除某个交易
func DeleteTransfer(mid int64) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	ss := fmt.Sprintf("delete from %s where id = ?", transferTable)

	_, err := db.Exec(ss, mid)
	if err != nil {
		fmt.Printf("DeleteTransfer err:%v", err)
	}

	return err
}