package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// CREATE TABLE `eth_transfer_key` (
// 	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
// 	`mto` varchar(64) NOT NULL UNIQUE COMMENT 'to地址',
//  `mfrom` varchar(64) NOT NULL COMMENT 'from',
// 	`raw` varchar(225) NOT NULL COMMENT '类型',
// 	`value` decimal(65,30) ZEROFILL DEFAULT '0' COMMENT '余额',
// 	`addTime` datetime DEFAULT NULL,
// 	PRIMARY KEY (`id`)
//   ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='eth地址表';

// SaveNewTransfer 保存新地址
func SaveNewTransfer(to, from, raw string, value int64) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	execStr := fmt.Sprintf("insert into %s(mto,mfrom,raw,value,addTime) values(?,?,?,?,Now())", transferTable)

	_, err := db.Exec(execStr, to, from, raw, value)
	if err != nil {
		log.Errorln("sql SaveNewTransfer err : ", err)
	}

	return err
}
