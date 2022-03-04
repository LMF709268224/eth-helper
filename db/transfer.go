package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// TODO 根据blocknumber排序

// CREATE TABLE `eth_transfer_key` (
// 	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
// 	`mto` varchar(64) NOT NULL UNIQUE COMMENT 'to地址',
//  `mfrom` varchar(64) NOT NULL COMMENT 'from地址',
// 	`txhash` varchar(225) NOT NULL COMMENT '交易哈希',
// 	`value` decimal(65,30) ZEROFILL DEFAULT '0' COMMENT '交易值',
// 	`blocknumber` bigint(20) ZEROFILL DEFAULT '0' COMMENT 'BlockNumber',
// 	PRIMARY KEY (`id`)
//   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='待验证交易表';

// SaveNewTransfer 保存待确认消息
func SaveNewTransfer(to, from, raw string, value int64, blocknumber uint64) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	execStr := fmt.Sprintf("insert into %s(mto,mfrom,raw,value,blocknumber) values(?,?,?,?,?)", transferTable)

	_, err := db.Exec(execStr, to, from, raw, value, blocknumber)
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
